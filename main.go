//
// bigmoney29
//
// Leonid Krashanoff
//
// ISC License
//

package main

import (
	"database/sql"
	"flag"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/time/rate"

	_ "github.com/mattn/go-sqlite3"
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
	e := echo.New()

	configPath := flag.String("c", "./config.toml", "`PATH` to config file to use.")
	flag.Parse()
	_, err := toml.DecodeFile(*configPath, &config)
	if err != nil {
		e.Logger.Fatal(err)
	}

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
	db, err := sql.Open("sqlite3", config.DBPath)
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

	// Serve our frontend files for any non-API GET request.
	e.GET("*", func(c echo.Context) error {
		return c.File("ui/build/index.html")
	})
	e.POST("/login", loginUser)

	// Userland API. Read assignments, grades, submit your work.
	api := e.Group("/largecurrency")
	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &UserClaim{},
		SigningKey: []byte(config.SigningKey),
	}))
	api.GET("/assignments", assignmentInformation)
	api.POST("/assignments/:id/submit", handleUpload)
	api.GET("/results/:id", serveResults)

	// Administrative API. Create, modify, delete users.
	admin := e.Group("/elongatedcash")
	admin.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(cc echo.Context) error {
			c := cc.(*Ctx)
			claims := c.Get("user").(*jwt.Token).Claims.(*UserClaim)
			if !claims.Admin {
				return c.NoContent(http.StatusUnauthorized)
			}
			return next(cc)
		}
	})
	api.POST("/user", createUser)
	api.DELETE("/user", deleteUser)

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
