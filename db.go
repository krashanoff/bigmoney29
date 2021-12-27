package main

import (
	"errors"

	bolt "go.etcd.io/bbolt"
)

var (
	metaBucket       = []byte("meta")
	assignmentBucket = []byte("assignment")
	userBucket       = []byte("user")
	buckets          = [][]byte{metaBucket, assignmentBucket, userBucket}
)

// Initialize the database.
func initDb(db *bolt.DB) error {
	// Create our buckets
	for _, bucketName := range buckets {
		if err := db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucket(bucketName)
			return err
		}); err != nil && err != bolt.ErrBucketExists {
			return err
		}
	}

	// Initialize the admin user, if this is our first-time setup.
	if err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(metaBucket)
		k, _ := b.Cursor().Seek([]byte("firstTimeComplete"))
		if k == nil {
			// Create the admin user.
			return nil
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Check if the user can be logged into.
func validateLogin(db *bolt.DB, username, password string) bool {
	if err := db.View(func(t *bolt.Tx) error {
		b := t.Bucket(userBucket).Bucket([]byte(username))
		if b == nil {
			return errors.New("user does not exist")
		}
		return nil
	}); err != nil {
		return false
	}
	return true
}
