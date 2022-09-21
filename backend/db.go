package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/blockloop/scan"
	"github.com/golang-jwt/jwt"

	_ "github.com/mattn/go-sqlite3"
)

var TABLEQUERIES = []string{
	`CREATE TABLE IF NOT EXISTS User (
		-- Username for the user
		username VARCHAR(255) UNIQUE NOT NULL,

		-- Hashed password of the user
		password VARCHAR(255) NOT NULL,

		-- Whether or not the user is an administrator
		admin BOOL NOT NULL,

		-- When the user was deleted (if at all)
		deleted TIMESTAMP,

		PRIMARY KEY (username)
	)`,
	`CREATE TABLE IF NOT EXISTS Assignment (
		-- ID of the assignment in the database
		id uuid_v4,

		-- Name of the assignment (to be used in submissions)
		name VARCHAR(255) UNIQUE NOT NULL,

		-- Points possible for the assignment
		points DOUBLE(5, 2),

		PRIMARY KEY (id)
	)`,
	`CREATE TABLE IF NOT EXISTS Submission (
		-- Unique submission ID
		id uuid_v4,
		
		-- ID for assignment
		assignment_id uuid_v4 NOT NULL,

		-- ID of submitter
		owner VARCHAR(255) NOT NULL,

		PRIMARY KEY (id),
		FOREIGN KEY (owner) REFERENCES User (username)
	)`,
	`CREATE TABLE IF NOT EXISTS CaseResult (
		-- Since this ID is determined by the grading script, it
		-- isn't our concern.
		case_id CHAR(10) NOT NULL,

		submission_id uuid_v4 NOT NULL,

		-- Points earned on this particular grading case.
		points_earned DOUBLE(5, 2),

		PRIMARY KEY (case_id, submission_id),
		FOREIGN KEY (submission_id) REFERENCES Submission (id)
	)`,
}

// Initialize the database schemas.
func initDb(db *sql.DB) error {
	for _, table := range TABLEQUERIES {
		if _, err := db.Exec(table); err != nil {
			return err
		}
	}
	return nil
}

// Confirm that the user with provided username and password does, in fact, exist in our
// database. If they do, return their user information.
func validateLogin(c *Ctx, username, password string) (*UserClaim, error) {
	tx, err := c.db.BeginTx(c.Request().Context(), &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	c.Logger().Infof("Username: %s Password: %s", username, password)

	rows, err := tx.QueryContext(c.Request().Context(), "SELECT password, admin FROM User WHERE ( username = ? )", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("incorrect credentials")
	}

	var (
		actualPassword string
		admin          bool
	)
	if err := rows.Scan(&actualPassword, &admin); err != nil || actualPassword != password {
		return nil, err
	}
	if rows.Next() {
		c.Logger().Warnf("Database is somehow corrupted for username '%s'", username)
		return nil, errors.New("database corrupted")
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.New("failed to commit")
	}

	claim := &UserClaim{
		admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
	}
	return claim, nil
}

// Add a user to the database.
func addUser(ctx context.Context, db *sql.DB, admin bool, uid, name, password string) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: false,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec("INSERT INTO User (username, password, admin) VALUES (?, ?, ?)", uid, name, admin); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// Delete a user from the database.
func removeUser(ctx context.Context, db *sql.DB, uid string) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: false,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec("UPDATE User SET deleted = CURRENT_TIMESTAMP WHERE username = '?'", uid); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

type Assignment struct {
	Name   string  `json:"name"`
	Points float64 `json:"points"`
}

// Get assignment information for the class.
func getAssignments(c *Ctx) ([]Assignment, error) {
	results, err := c.db.Query("SELECT * FROM Assignment")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	assignments := []Assignment{}
	if !results.Next() {
		return assignments, nil
	}
	if err := scan.Rows(&assignments, results); err != nil {
		return nil, err
	}

	return assignments, nil
}

func addAssignment(ctx context.Context, db *sql.DB, name string, totalPoints float64) error {
	if _, err := db.ExecContext(ctx, "INSERT INTO Assignment (name, points) VALUES (?, ?)", name, totalPoints); err != nil {
		return err
	}
	return nil
}

func removeAssignment(ctx context.Context, db *sql.DB, name string) error {
	if _, err := db.ExecContext(ctx, "DELETE FROM Assignment WHERE name = '?'", name); err != nil {
		return err
	}
	return nil
}
