# bigmoney29

tremendous dub bossman

![snorlax sitting](docs/snorlax.png)

## About

`bigmoney29` is an automated grading solution for small programming classes. It is
reasonably small and easy to use.

## Configuration/Setup

The server is distributed as a static binary and webpacked React application. To
start it, simply run the binary.

Configure the server through `config.toml`. The [example file](config.toml.example)
is heavily commented. All options should be set. Changes to the configuration file
are reflected on server restart.

## Admin Interface

You can run administrative operations on the server through its SSH interface.
All operations are conducted by sending a file over SSH. Command names are
case insensitive, and double quotes can be replaced by single quotes where
seen. Mixing double and single quotes within a single string is undefined
behavior (e.g., `"Student's Name"`).

### Add Students

```sh
echo '+student UID "Student Name" TemporaryPassword' > mystudents
echo '+student UID2 "Other Student" TempPassword' >> mystudents
ssh -p 8082 localhost <mystudents
```

### Remove Students

```sh
ssh -p 8082 localhost '-student UID2'
```

### Add Assignment

```sh
ssh -p 8082 localhost '+assignment "assignment name"'
```

### Modify a Student's Final Score

```sh
ssh -p 8082 localhost '+score UID2 "ASSIGNMENT NAME" newscore'
```

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
* A grading script must, for each test case:
  * Print the name of the test case.
  * Print supplementary help messages to stderr.
  * Terminate the test case by performing the following:
    * Print the weight for the test case and the program's score for the test case as
      a 64-bit floating point number.
    * Print an extra newline.
* The final score is calculated as the weighted average of these test cases.

### Example

Let's assume that we have one assignment in our class named "Assignment One". Let's also
assume that we wrote a grading script for assignment one that will give 100 points to any
student who submitted a single file and 50 points to all students regardless. This script
takes a single argument: the directory to check.

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

Then, we can write a driver script for this:

```sh
#!/bin/sh
case "$1" in
  "Assignment One")
    ./grade_assignment_one.sh $2
    ;;
  *)
    echo "Called with an unknown assignment." 1>&2
    ;;
esac
```

Finally, we configure the server to use our driver script.

```toml
grading_script = "./driver.sh"
```

## License

ISC License. 2021, Leonid Krashanoff.

## What is the name

I made it up.
