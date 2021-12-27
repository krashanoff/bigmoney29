# bigmoney29

tremendous dub bossman

![snorlax sitting](docs/snorlax.png)

## About

`bigmoney29` is an automated grading solution for small programming classes. It is
reasonably small and easy to maintain.

## TODO

* User association
* Use ICO instead of PNG for favicon
* Alternative submission formats
* Grade distributions + data analysis views

## Configuration/Setup

Configure the server through `config.toml`. The [example file](config.toml.example)
is heavily commented. All options should be set.

## Admin Accounts

You can issue admin accounts through the superuser created when the server is first started.
The username and password are both "admin".

`bigmoney` hosts an administrative API on a port of your choosing to get final
submissions.
* Get final submissions

## Grading Scripts

Grading scripts are highly permissive by design. You can do whatever you really
want to in these... and so can students. Exercise caution.

`bigmoney` provides the following guarantees to your grading script:

* It is guaranteed that the name of the assignment a student is submitting to
  (e.g., "Homework 1") will be provided as the first argument to the script
  (`$1`) or program (`argv[1]`).
* It is guaranteed that the path to a student's submitted code will be provided as the
  second argument to the script (`$2`) or program (`argv[2]`). The path provided will
  have chmod `0640` by default, and be created by the user `bigmoney` is running under.

In exchange, your grading script must provide the following functionality:
* A grading script **must** have chmod `0700` for the `bigmoney` user.
* A grading script must handle the following option `--summary`.
  * Output each grading scheme.
* A grading script must, for each test case:
  * Print supplementary help messages to stderr.
  * Print the weight for the test case and the program's score for the test case as
    a 64-bit floating point number.
* The final score is calculated as the weighted average of these test cases.

For example, a grading script that will give 100 points to any student who submitted a
single file and 50 points to all students regardless will look like:

```sh
#!/bin/sh
cd $1
echo "0.5 1"
if [ "$(ls -a | wc -l)" == "4" ] then
  echo "Single file present" 1>&2
  echo "0.5 1"
else
  echo "Missing one file" 1>&2
  echo "0.5 0"
fi
```

A grading program that performs the same operation might look like:

```go
package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "os"
)

func main() {
  flag.Parse()
  assignmentName, chDir := flag.Arg(0), flag.Arg(1)
  if assignmentName == "" || chDir == "" {
    fmt.Println("1 0")
    return
  }

  files, err := ioutil.ReadDir(chDir)
  if err != nil {
    fmt.Fprintln(os.Stderr, "Failed to read the directory for some reason. Email your professor.")
    fmt.Println("1 0")
  }
  fmt.Println("0.5 1")
  if len(files) == 1 {
    fmt.Fprintln(os.Stderr, "Single file present")
    fmt.Println("0.5 1")
  } else {
    fmt.Fprintln(os.Stderr, "Missing one file")
    fmt.Println("0.5 0")
  }
}
```

## License

ISC License. 2021, Leonid Krashanoff.

## What is the name

I made it up.
