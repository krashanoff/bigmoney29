package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/gliderlabs/ssh"
	"github.com/labstack/echo/v4"
)

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
			if strings.HasPrefix(line, "#") {
				continue
			}

			components := strings.Split(line, " ")
			commandName, operands := components[0], components[1:]
			logger.Infof("Command name: '%s'. Operands: '%v'", commandName, operands)

			switch strings.ToLower(commandName) {
			case "+student":
				if len(operands) != 3 {
					failed++
					continue
				}
				if err := addUser(s.Context(), db, false, operands[0], operands[1], operands[2]); err != nil {
					failed++
					continue
				}
				lineWriter.WriteString(fmt.Sprintf("Added student '%s'\n", "name"))
			case "-student":
				if len(operands) != 1 {
					failed++
					continue
				}
				if err := removeUser(s.Context(), db, operands[0]); err != nil {
					failed++
					continue
				}
				lineWriter.WriteString(fmt.Sprintf("Removed student '%s'\n", "name"))
			case "+assignment":
				if len(operands) != 2 {
					failed++
					continue
				}
				totalScore, _ := strconv.ParseFloat(operands[1], 64)
				if err := addAssignment(s.Context(), db, operands[0], totalScore); err != nil {
					failed++
					continue
				}
				lineWriter.WriteString(fmt.Sprintf("Created assignment '%s'\n", "name"))
			case "-assignment":
				if len(operands) != 1 {
					failed++
					continue
				}
				if err := removeAssignment(s.Context(), db, operands[0]); err != nil {
					failed++
					continue
				}
				lineWriter.WriteString(fmt.Sprintf("Removed assignment '%s'\n", "name"))
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
