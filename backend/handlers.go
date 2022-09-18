package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type UserClaim struct {
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

type Ctx struct {
	echo.Context
	db       *sql.DB
	jobQueue chan<- Job
}

// Log the user in, providing a JWT.
func loginUser(cc echo.Context) error {
	c := cc.(*Ctx)

	// Login page allowing use of the bigmoney API.
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Validate the username and password against our database.
	claim, err := validateLogin(c, username, password)
	if err != nil {
		c.Logger().Warnf("User '%s' could not login: %v", username, err)
		return c.NoContent(http.StatusForbidden)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := token.SignedString([]byte(config.SigningKey))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// Get assignments for class
func assignmentInformation(cc echo.Context) error {
	c := cc.(*Ctx)
	c.Logger().Info("Checking assignments")
	assn, err := getAssignments(c)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, assn)
}

// Handle an upload to the server. Validate the upload, generate a new Job,
// and proceed.
func handleUpload(cc echo.Context) error {
	c := cc.(*Ctx)

	assn := c.FormValue("assignment")
	if assn == "" {
		return c.String(http.StatusBadRequest, "Invalid assignment")
	}
	file, err := c.FormFile("file")
	if err != nil || file == nil || file.Size == 0 {
		return c.String(http.StatusBadRequest, "Invalid file field")
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
			// Place run, then send over websocket.
			fmt.Print(id)
		}
		tick.Stop()
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
