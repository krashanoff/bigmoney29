package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	bolt "go.etcd.io/bbolt"
)

// Run the grading script for a single Job.
func gradingScript(db *bolt.DB, job Job, results chan<- Result, occupied <-chan bool) error {
	defer func() { <-occupied }()
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
	dirPfx := "./run"
	dir := path.Join(dirPfx, dirPfx+"-"+strings.ReplaceAll(job.id, "-", ""))
	if err := os.MkdirAll(dir, 0640); err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// Copy all of our files from the gzip.
	for header, err := tr.Next(); err == nil; header, err = tr.Next() {
		tmpFile, err := os.Create(path.Join(dir, header.Name))
		os.Chmod(tmpFile.Name(), 0640)
		if err != nil {
			fmt.Println(err)
		}
		if _, err := io.Copy(tmpFile, tr); err != nil {
			fmt.Println(err)
		}
	}

	output, err := exec.Command(config.GradingScript, dir).Output()
	if err != nil {
		return err
	}
	fmt.Print(string(output))
	return nil
}
