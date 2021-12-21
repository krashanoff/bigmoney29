# bigmoney29

tremendous dub bossman

![snorlax sitting](docs/snorlax.png)

Grading server for small programming classes.

## TODO

* User accounts + admin accounts
* Grade distributions + data analysis views

## Configuration/Setup

Configure the server through `config.toml`. The [example file](config.toml.example)
is heavily commented.

## Grading Scripts

Grading scripts are highly permissive by design. You can do whatever you really
want to in these... and so can students! Exercise caution. Maybe use a docker
image -- all up to you!

* It is guaranteed that the path to a student's submitted code will be provided as the
  first argument to the script (`$1`) or program (`argv[1]`).
* For each test case in a grading script:
  * Print supplementary help messages to stderr.
  * Print the weight for the test case and the program's score for the test case on
    the range 0.0-100.0.
* The final score is calculated as the weighted average of these test cases.

For example, a grading script that will give 100 points to any student who submitted a
single file and 50 points to all students regardless will look like:

```sh
#!/bin/sh
cd $1
echo "0.5 100"
if [ "$(ls -a | wc -l)" == "4" ] then
  echo "0.5 100"
else
  echo "0.5 0"
fi
```

## License

ISC License. 2021, Leonid Krashanoff.
