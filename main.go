//
// athome
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
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/net/websocket"
)

// Maximum number of parallel job tasks.
const MAXJOBS = 5

// Jobs queued for execution.
var jobQueue chan Job

var db *bolt.DB

type Job struct {
	id         string
	file       io.Reader
	assignment string
}

type Result struct {
	Fail bool `json:"fail"`
}

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.GET("*", func(c echo.Context) error {
		return c.File("ui/build/index.html")
	})
	e.POST("/upload", func(c echo.Context) error {
		name := c.FormValue("name")
		assn := c.FormValue("assignment")
		if name == "" || assn == "" {
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
		jobQueue <- Job{
			id:         requestId,
			file:       fileData,
			assignment: assn,
		}
		return c.String(http.StatusCreated, requestId)
	})
	e.GET("/results/:id", getResults)

	handle, err := bolt.Open("tests.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	if err := initDb(handle); err != nil {
		e.Logger.Error(err)
		return
	}
	db = handle

	// Spawn task runner
	jobQueue = make(chan Job, MAXJOBS)
	go func() {
		for job := range jobQueue {
			// TODO: limit number of jobs that can run in parallel
			go func(job Job) {
				zr, err := gzip.NewReader(job.file)
				if err != nil {
					e.Logger.Error(err)
					db.Update(func(t *bolt.Tx) error {
						data, err := json.Marshal(Result{Fail: true})
						if err != nil {
							return err
						}
						t.Bucket([]byte("runs")).Put([]byte(job.id), data)
						return nil
					})
					return
				}
				defer zr.Close()
				tr := tar.NewReader(zr)

				// Prepare our temporary directory for running files.
				dir := os.TempDir()
				e.Logger.Debugf("Temp dir is: %v", dir)
				for header, err := tr.Next(); err == nil; header, err = tr.Next() {
					tmpFile, err := os.CreateTemp(dir, header.Name)
					if err != nil {
						break
					}
					io.Copy(tmpFile, tr)
				}

				// Once a temporary directory is established, we add it to our DB.
				db.Update(func(tx *bolt.Tx) error {
					return nil
				})

				cmd := exec.Command("cargo", "test", "--manifest-path", path.Join(dir, "Cargo.toml"))
				stdout, _ := cmd.StdoutPipe()
				if err := cmd.Run(); err != nil {
					e.Logger.Errorf("%v\n", err)
				}
				fmt.Println(job, stdout)
			}(job)
		}
	}()

	// Serve
	e.Logger.Fatal(e.Start(":8080"))
}

func getResults(c echo.Context) error {
	id := c.Param("id")
	c.Logger().Debug("Starting websocket")
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		// Listen in on our tests, displaying output as it arrives.
		tick := time.NewTicker(500 * time.Millisecond)
		for range tick.C {
			db.View(func(t *bolt.Tx) error {
				c.Logger().Debug("starting")
				b := t.Bucket([]byte("runs"))
				data := b.Get([]byte(id))

				c.Logger().Debugf("Checked and got %s", string(data))

				err := websocket.Message.Send(ws, data)
				if err != nil {
					c.Logger().Error(err)
				}
				return nil
			})
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
