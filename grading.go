package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// A job to run the grading script on.
type Job struct {
	id         string
	file       io.Reader
	assignment string
}

// Result from running a job.
type Result struct {
	TestName string  `json:"testName"`
	Score    float64 `json:"score"`
	Pass     bool    `json:"pass"`
	Msg      string  `json:"msg"`
}

// Run the grading script for a single Job.
func gradingScript(db *sql.DB, job Job, results chan<- Result, occupied <-chan bool) error {
	defer func() { <-occupied }()
	zr, err := gzip.NewReader(job.file)
	if err != nil {
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

// Parse lines from a reader into Results, returning an error
// if output was malformed.
func parseOutput(stdout io.Reader) ([]Result, error) {
	stdoutLines := bufio.NewScanner(stdout)
	results := make([]Result, 0)
	for stdoutLines.Scan() {
		resultScore := float64(0)
		name := stdoutLines.Text()
		lines := make([]string, 0)
		for stdoutLines.Scan() {
			text := stdoutLines.Text()
			if text == "" {
				break
			}
			lines = append(lines, text)
		}

		scoreLine := lines[len(lines)-1]
		lines = lines[:len(lines)-1]
		scoreParts := strings.Split(scoreLine, " ")

		weight, err := strconv.ParseFloat(scoreParts[0], 64)
		if err != nil {
			return nil, err
		}
		score, err := strconv.ParseFloat(scoreParts[1], 64)
		if err != nil {
			return nil, err
		}
		resultScore += weight * score
		results = append(results, Result{
			TestName: name,
			Msg:      strings.Join(lines, "\n"),
			Score:    resultScore,
		})
	}
	return results, nil
}
