package main

import bolt "go.etcd.io/bbolt"

// Initialize the database.
func initDb(db *bolt.DB) error {
	if err := db.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucket([]byte("runs"))
		return err
	}); err != nil && err != bolt.ErrBucketExists {
		return err
	}
	return nil
}
