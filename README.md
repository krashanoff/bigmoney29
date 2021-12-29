# bigmoney29

tremendous dub bossman

![snorlax sitting](docs/snorlax.png)

## About

`bigmoney29` is an automated grading solution for small programming classes. It is
reasonably small and easy to maintain.

## TODO

* Hash passwords (lol)
* Use ICO instead of PNG for favicon (lol)
* Alternative submission formats

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
  * Print the name of the test case.
  * Print supplementary help messages to stderr.
  * Terminate the test case by performing the following:
    * Print the weight for the test case and the program's score for the test case as
      a 64-bit floating point number.
    * Print an extra newline.
* The final score is calculated as the weighted average of these test cases.

For example, a grading script that will give 100 points to any student who submitted a
single file and 50 points to all students regardless will look like:

```sh
#!/bin/sh
cd $1
echo "Submitted"
echo "0.5 1\n"

echo "Contains one file"
if [ "$(ls -a | wc -l)" == "4" ]; then
  echo "0.5 1\n"
else
  echo "Missing one file"
  echo "0.5 0\n"
fi
```

## License

ISC License. 2021, Leonid Krashanoff.

## What is the name

I made it up.
