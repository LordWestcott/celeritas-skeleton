mod file is initialised to 
`github.com/lordwestcott/celeritas`
This is because we wish people to pull it.

Using github.com/joho/godotenv to read what is in the .env file. Used in Celeritas.New.

Using github.com/go-chi/chi to implement routing.
`go get -u github.com/go-chi/chi/v5`
And for middleware -> 
`go get github.com/go-chi/chi/v5/middleware`

Something new learnt:
if you have something that is of type `interface{}`
you can cast it as such:
(in this example `data` is of type `interface{}` and we are casting it to `*TemplateData`)

`td = data.(*TemplateData)` //This casts the interface{} to *TemplateData

For Jet renderering we are using: 
`go get github.com/CloudyKit/jet/v6`

OH MY GOD PLEASE TAKE NOTE OF THIS:
`go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out`
Coverage shown in browser.

For Sessions we are using
`go get github.com/alexedwards/scs/v2`
Later for Session Stores with postgres we are using:
`go get github.com/alexedwards/scs/postgresstore`
Same for MySQL:
`github.com/alexedwards/scs/mysqlstore`
For Redis Session Store its:
`go get github.com/alexedwards/scs/redisstore`

For Database (Postgres) we are using:
`go get github.com/jackc/pgconn`
`go get github.com/jackc/pgx/v4`
`go get github.com/jackc/pgx/v4/stdlib`

For Coloring the Terminal, we are using the package:
`go get github.com/fatih/color`

For migrations we will be using:
`go get github.com/golang-migrate/migrate/v4`
`go get github.com/golang-migrate/migrate/v4/database/mysql`
`go get github.com/golang-migrate/migrate/v4/database/postgres`
`go get github.com/golang-migrate/migrate/v4/source/file`

Tool to turn strings into different case patterns:
`go get github.com/iancoleman/strcase`

For help turning words into their plural counterparts:
`go get github.com/gertd/go-pluralize`

For validation we are using govalidator:
`go get github.com/asaskevich/govalidator`

For Redis Caching (and interacting with redis)
we are using the redigo package:
`go get github.com/gomodule/redigo/redis`

For Badger Caching (and interacting with BadgerDB) we are using:
`go get github.com/dgraph-io/badger/v3`

For Redis Unit testing we are using a cool little package that automatically spins up a mini redis server specifically for unit testing.
This means we don't need to spin up docker images like out postgres integration testing.
`go get github.com/alicebob/miniredis/v2`

For CSRF Token protection implementation we are using this package:
`go get github.com/justinas/nosurf`

For CRON jobs - (specifically BadgerDB cleanup)
We are using robfig's cron package:
`go get github.com/robfig/cron/v3`

For Emails (APIs)
the course is using this:
`go get github.com/ainsleyclark/go-mail@v1.0.3`
but it has since been updated to this:
`go get github.com/ainsleyclark/go-mail`
Other email libs we are going to use:
- go-simple-mail `go get github.com/xhit/go-simple-mail/v2`
    - Easy way to send mail.
    - Creates a mail client, add attahments etc
- go-premailer `go get github.com/vanng822/go-premailer/premailer`
    - Takes any styling on a html email message and inlines the css.

For mail testing you can use something like [mailtrap](https://mailtrap.io/)
we are going to do it locally with [mailhog](https://github.com/mailhog/MailHog)

We are also pulling a docker image to test mail functionality using the mailhog image.
We have already implemented docker integration testing in the mayapp project. But just to put it here too, we use:
`go get github.com/ory/dockertest/v3`
`go get github.com/ory/dockertest/v3/docker`

For urlsigning we are using:
`go get github.com/bwmarrin/go-alone`