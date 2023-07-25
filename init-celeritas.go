package main

import (
	"log"
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"
	"os"

	"github.com/lordwestcott/celeritas"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//init celeritas
	cel := &celeritas.Celeritas{}
	err = cel.New(path)
	if err != nil {
		log.Fatal(err)
	}

	cel.AppName = "myapp"

	middlewares := &middleware.Middleware{
		App: cel,
	}

	handlers := &handlers.Handlers{
		App: cel,
	}

	cel.InfoLog.Println("Debug is set to: ", cel.Debug)

	app := &application{
		App:        cel,
		Handlers:   handlers,
		Middleware: middlewares,
	}

	app.App.Routes = app.routes() //Adding routes to existing routes.

	app.Models = data.New(app.App.DB.Pool)
	app.Middleware.Models = app.Models
	handlers.Models = app.Models

	return app

}
