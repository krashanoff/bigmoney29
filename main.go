//
// bigmoney29
//
// Leonid Krashanoff
//
// ISC License
//

package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/net/websocket"
)

var config Config

type Config struct {
	GradingScript      string `toml:"grading_script"`
	Address            string `toml:"address"`
	ResultsRefreshRate uint   `toml:"results_refresh_rate"`
	MaxJobs            uint   `toml:"max_jobs"`
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
	e.Use(middleware.Gzip())
	e.Use(middleware.BodyLimit("10K"))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Apply our context.
	jobQueue := make(chan Job, config.MaxJobs)
	results := make(chan Result, config.MaxJobs)
	db, err := bolt.Open("athome.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := initDb(db); err != nil {
		e.Logger.Error(err)
		return
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

	e.GET("*", func(c echo.Context) error {
		return c.File("ui/build/index.html")
	})
	e.POST("/upload", handleUpload)
	e.GET("/results/:id", serveResults)

	// Spawn our task runner
	go func() {
		for job := range jobQueue {
			go gradingScript(db, job, results)
		}
	}()

	// Serve
	e.Logger.Fatal(e.Start(config.Address))
}

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

func gradingScript(db *bolt.DB, job Job, results chan<- Result) error {
	fmt.Println("Starting grading script")
	zr, err := gzip.NewReader(job.file)
	if err != nil {
		db.Update(func(t *bolt.Tx) error {
			data, err := json.Marshal([]Result{{Fail: true}})
			if err != nil {
				return err
			}
			t.Bucket([]byte("runs")).Put([]byte(job.id), data)
			return nil
		})
		return err
	}
	defer zr.Close()
	tr := tar.NewReader(zr)

	// Prepare our temporary directory for running files.
	dir := os.TempDir()
	os.Mkdir(dir, 0640)
	for header, err := tr.Next(); err == nil; header, err = tr.Next() {
		fmt.Printf("Checking file %s\n", header.Name)
		tmpFile, err := os.CreateTemp(dir, header.Name)
		if err != nil {
			fmt.Println("failed, breaking", err)
			break
		}
		io.Copy(tmpFile, tr)
	}
	fmt.Println("dir is", dir)

	// Once a temporary directory is established, we register it to our DB.
	// db.Update(func(tx *bolt.Tx) error {
	// 	return nil
	// })

	fmt.Println("path is", config.GradingScript)
	cmd := exec.Command(config.GradingScript, dir)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Print(string(output))
	return nil
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

func initDb(db *bolt.DB) error {
	if err := db.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucket([]byte("runs"))
		return err
	}); err != nil && err != bolt.ErrBucketExists {
		return err
	}
	return nil
}
