package main

import (
	"database/sql"
	"errors"
	"time"

	"github.com/blockloop/scan"
	"github.com/golang-jwt/jwt"
)

var TABLEQUERIES = []string{
	`CREATE TABLE User (
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		admin BOOL NOT NULL,
		deleted TIMESTAMP,
		PRIMARY KEY (username)
	)`,
	`CREATE TABLE Assignment (
		-- Name of the assignment (to be used in submissions)
		name VARCHAR(255) UNIQUE NOT NULL,

		-- Points possible for the assignment
		points DOUBLE(5, 2),
		PRIMARY KEY (name)
	)`,
	`CREATE TABLE Submission (
		-- Unique submission ID
		id CHAR(10),
		
		-- ID for assignment
		assignment VARCHAR(50) NOT NULL,

		-- ID of submitter
		owner VARCHAR(255) NOT NULL,

		-- Total points earned in this submission
		points_earned DOUBLE(5, 2),
		PRIMARY KEY (id)
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

	rows, err := tx.QueryContext(c.Request().Context(), "SELECT ( password, admin ) FROM User WHERE ( username = ? )", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
func addUser(c *Ctx, admin bool, uid, name, password string) error {
	tx, err := c.db.BeginTx(c.Request().Context(), &sql.TxOptions{
		ReadOnly: false,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec("INSERT INTO User (?, ?, ?, ?)", admin, uid, name, password); err != nil {
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
	tx, err := c.db.BeginTx(c.Request().Context(), &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	results, err := tx.Query("SELECT * FROM Assignment")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var assignments []Assignment
	if err := scan.Rows(&assignments, results); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return assignments, nil
}
