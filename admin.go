package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/gliderlabs/ssh"
	"github.com/labstack/echo/v4"
)

var matchCommand = regexp.MustCompile(`([\+\-](student|assignment|score)) (.*) (".*"|'.*') (.*)`)

func startAdmin(e *echo.Echo, db *sql.DB) (*ssh.Server, error) {
	s, err := wish.NewServer(wish.WithAddress("localhost:8082"), wish.WithPasswordAuth(checkAdmin), wish.WithMiddleware(
		bm.Middleware(processCommandsMiddleware(db)),
	))
	if err != nil {
		return nil, err
	}
	return s, nil
}

func checkAdmin(ctx ssh.Context, password string) bool {
	return password == config.Password
}

func processCommandsMiddleware(db *sql.DB) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		lines := bufio.NewScanner(s)
		lineWriter := bufio.NewWriter(s)
		successful, failed := 0, 0
		for lines.Scan() {
			line := lines.Text()
			submatch := matchCommand.FindStringSubmatch(line)
			if submatch == nil {
				lineWriter.WriteString("Invalid command: " + line)
				failed++
				continue
			}

			switch strings.ToLower(submatch[1]) {
			case "+student":
				// TODO: add user to class
				if err := addUser(s.Context(), db, false, "22", "leo", "pw"); err != nil {
					failed++
					fmt.Println("Failed", err)
					continue
				}
				fmt.Println("Processed")
			case "-student":
			case "+assignment":
			case "-assignment":
			case "+score":
			case "-score":
			}
			successful++
		}
		lineWriter.WriteString(fmt.Sprintf("OK: %d FAILED: %d\n", successful, failed))
		return nil, nil
	}
}
