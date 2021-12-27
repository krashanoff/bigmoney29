//
// bigmoney29
//
// Leonid Krashanoff
//
// ISC License
//

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/net/websocket"
	"golang.org/x/time/rate"
)

var config Config

// See config.toml.example for definitions of this struct's fields.
type Config struct {
	GradingScript      string `toml:"grading_script"`
	Address            string `toml:"address"`
	BodyLimit          string `toml:"body_limit"`
	DBPath             string `toml:"db_path"`
	ResultsRefreshRate uint   `toml:"results_refresh_rate"`
	MaxJobs            uint   `toml:"max_jobs"`
	RateLimit          uint   `toml:"rate_limit"`
	SigningKey         string `toml:"signing_key"`
}

func init() {
	config = Config{
		Address:            ":8080",
		BodyLimit:          "10K",
		DBPath:             "bigmoney.db",
		ResultsRefreshRate: 3000,
		MaxJobs:            5,
		RateLimit:          5,
	}
}

func main() {
	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Logger())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(config.RateLimit))))
	e.Use(middleware.Gzip())
	e.Use(middleware.BodyLimit(config.BodyLimit))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// Apply our custom context to each request.
	jobQueue := make(chan Job, config.MaxJobs)
	results := make(chan Result, config.MaxJobs)
	db, err := bolt.Open(config.DBPath, 0600, nil)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()
	if err := initDb(db); err != nil {
		e.Logger.Fatal(err)
	}
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&Ctx{
				c,
				db,
				jobQueue,
			})
		}
	})

	// Serve our JavaScript files for any GET request.
	e.GET("*", func(c echo.Context) error {
		return c.File("ui/build/index.html")
	})

	e.POST("/upload", handleUpload)
	e.GET("/ws/:id", serveResults)
	e.POST("/cash/login", func(c echo.Context) error {
		// Login page allowing use of the bigmoney API.
		username := c.FormValue("username")
		password := c.FormValue("password")

		// Validate the username and password against our database.
		if !validateLogin(db, username, password) {
			return c.NoContent(http.StatusUnauthorized)
		}

		// TODO: get user information from database
		claim := &UserClaim{
			true,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		t, err := token.SignedString([]byte(config.SigningKey))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
	})

	api := e.Group("/cash")
	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &UserClaim{},
		SigningKey: []byte(config.SigningKey),
	}))
	api.GET("/assignments", func(c echo.Context) error {
		// TODO: get assignments for class
		return c.String(http.StatusOK, "{\"assignment\": \"\"}")
	})
	api.GET("/assignments/:assignmentID", func(c echo.Context) error {
		// TODO: get score for assignment
		return c.String(http.StatusOK, "Welcome")
	})
	api.PUT("/user", func(c echo.Context) error {
		return c.String(http.StatusOK, "Test")
	})
	api.POST("/user", func(c echo.Context) error {
		claims := c.Get("user").(*jwt.Token).Claims.(*UserClaim)
		if !claims.Admin {
			return c.NoContent(http.StatusUnauthorized)
		}

		// TODO: create user(s)
		var body []struct {
			Admin    bool   `json:"admin"`
			Name     string `json:"name"`
			UID      string `json:"UID"`
			Password string `json:"password"`
		}
		c.Bind(&body)
		return nil
	})
	api.DELETE("/user", func(c echo.Context) error {
		claims := c.Get("user").(*jwt.Token).Claims.(*UserClaim)
		if !claims.Admin {
			return c.NoContent(http.StatusUnauthorized)
		}

		var body []struct {
			Name string `json:"name"`
			UID  string `json:"UID"`
		}
		c.Bind(&body)

		return c.NoContent(http.StatusOK)
	})

	// Spawn our task runner
	go func() {
		occupied := make(chan bool, config.MaxJobs)
		for job := range jobQueue {
			occupied <- true
			go gradingScript(db, job, results, occupied)
		}
	}()

	// Serve
	e.Logger.Fatal(e.Start(config.Address))
}

type UserClaim struct {
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

// Handle an upload to the server. Validate the upload, generate a new Job,
// and proceed.
func handleUpload(cc echo.Context) error {
	c := cc.(*Ctx)

	// name := c.FormValue("name")
	assn := c.FormValue("assignment")
	if assn == "" {
		return c.String(http.StatusBadRequest, "Invalid assignment number")
	}
	file, err := c.FormFile("file")
	if err != nil || file == nil || file.Size == 0 {
		return c.String(http.StatusBadRequest, "Failed parsing file field")
	}

	fileData, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read .tar.gz")
	}

	// Generate the job.
	requestId := uuid.NewString()
	c.jobQueue <- Job{
		id:         requestId,
		file:       fileData,
		assignment: assn,
	}
	return c.String(http.StatusCreated, requestId)
}

type Ctx struct {
	echo.Context
	db       *bolt.DB
	jobQueue chan<- Job
}

type Job struct {
	id         string
	file       io.Reader
	assignment string
}

type Result struct {
	TestName string `json:"testName"`
	Score    uint64 `json:"score"`
	Fail     bool   `json:"fail"`
	Msg      string `json:"msg"`
}

// Push results to a client at regular intervals. If the results are already
// available, then just send them the results and terminate.
func serveResults(cc echo.Context) error {
	c := cc.(*Ctx)

	id := c.Param("id")
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		// Listen in on our tests, displaying output as it arrives.
		tick := time.NewTicker(time.Duration(config.ResultsRefreshRate) * time.Millisecond)
		for range tick.C {
			err := c.db.View(func(t *bolt.Tx) error {
				b := t.Bucket([]byte("runs"))
				data := b.Get([]byte(id))

				var results []Result
				json.Unmarshal(data, &results)

				err := websocket.JSON.Send(ws, results)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				c.Logger().Error(err)
				break
			}
		}
		tick.Stop()
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
