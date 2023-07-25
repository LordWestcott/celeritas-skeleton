in go.mod ->
replace github.com/lordwestcott/celeritas => ../celeritas

This means that whenever we reference github.com/lordwestcott/celeritas in code, it will actually pull from the sibling celeritas directory.


go mod vendor
This creates vendor folder - Keeps a local folder of everything 3rd party that we are building.

We can utilize this to sync this with package changes.
Check out the MakeFile, now we just have to 'make run'.

For ORM we are using Upper:
Postgres:
`go get github.com/upper/db/v4/adapter/postgresql`
MySQL:
`got get github.com/upper/db/v4/adapter/mysql`

for mocking SQL
`go get github.com/DATA-DOG/go-sqlmock`

Integration testing:
Working with docker images, we will use the package:
`go get github.com/ory/dockertest/v3`
`go get github.com/ory/dockertest/v3/docker`

