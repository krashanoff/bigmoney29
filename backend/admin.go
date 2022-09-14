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
	s, err := wish.NewServer(wish.WithAddress(config.AdminAddress), wish.WithHostKeyPath(config.PrivateKeyPath), wish.WithPasswordAuth(checkAdmin), wish.WithMiddleware(
		bm.Middleware(processCommandsMiddleware(e.Logger, db)),
	))
	if err != nil {
		return nil, err
	}
	return s, nil
}

func checkAdmin(ctx ssh.Context, password string) bool {
	return password == config.Password
}

func processCommandsMiddleware(logger echo.Logger, db *sql.DB) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		lines := bufio.NewScanner(s)
		lineWriter := bufio.NewWriter(s)
		logger.Infof("Successful admin login from %v using username %s", s.RemoteAddr(), s.User())

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
				if err := addUser(s.Context(), db, false, "22", "leo", "pw"); err != nil {
					failed++
					continue
				}
				lineWriter.WriteString(fmt.Sprintf("Added student '%s'\n", "name"))
			case "-student":
				if err := removeUser(s.Context(), db, "22"); err != nil {
					failed++
					continue
				}
				lineWriter.WriteString(fmt.Sprintf("Removed student '%s'\n", "name"))
			case "+assignment":
				if err := addAssignment(s.Context(), db, "name", 0.0); err != nil {
					failed++
					continue
				}
				lineWriter.WriteString(fmt.Sprintf("Created assignment '%s'\n", "name"))
			case "-assignment":
				if err := removeAssignment(s.Context(), db, "name"); err != nil {
					failed++
					continue
				}
				lineWriter.WriteString(fmt.Sprintf("Removed assignment '%s'\n", "name"))
			case "+score":
				if err := addScore(s.Context(), db, "name"); err != nil {
					failed++
					continue
				}
			default:
				lineWriter.WriteString(fmt.Sprintf("WARNING: Unknown transaction '%s'\n", line))
			}
			successful++
			lineWriter.Flush()
		}
		lineWriter.WriteString(fmt.Sprintf("Transaction Summary:\nOK: %d FAILED: %d\n", successful, failed))
		lineWriter.Flush()
		return nil, nil
	}
}
