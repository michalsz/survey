# Setting up your environment in under 5 minutes
Make sure you have docker and docker-compose installed on your computer.
Then, run `docker-compose run --service-ports app /bin/bash`.

You'll be presented with a command line with all tools already preinstalled,
connected to a test database.

Let's start by setting up the database schema by running `./setup-db`. You can
run `./db-connect` to access the postgres db console.

Once we're done, we can install the dependencies by running `dep ensure -v`

Congratulations! You're now ready to start developing. Yay!

# Meeting the code
The codebase is a stripped down project with only the things you need for
developing the feature at hand.

The code is structured into `controllers/`, `models/` and `serializers/`
directories. Controllers handle the request, serializers are used to convert
models to/from what gets transfered over the network and models deal with the
database. Framework called BeeGo is used, although not much of your code will
touch it.

Testing is of huge importance and lives in the `unittest/` directory.
The `init.go` file gets executed everytime the test package gets loaded and
sets up BeeGo.

The recommended way of running tests is `go test -v ./unittest`. If you have
followed all of the previous steps correctly, the test should pass. Give it a
try!

If you need to run a specific test, then you can run
`go test -v ./test -run ^TestName$`. Note that `^TestName$` is a regular
expression - you can just write `TestName`, but in such case other tests with
similar names may get executed as well.

To run the code, you can just run `go run main.go` and then access it at
`http://localhost:8080`.

# Let's get down to work!
Now you should be ready to complete the work at hand. Please check the file
`TASK.md` for details.