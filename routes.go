package main

import (
	"fmt"
	"myapp/data"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lordwestcott/celeritas/mailer"
)

// This basically Adds routes to the existing routes.
func (a *application) routes() *chi.Mux {
	// Middleware must come before any routes
	a.use(a.Middleware.CheckRemember)

	// Add routes here
	a.get("/", a.Handlers.Home)
	a.get("/go-page", a.Handlers.GoPage)
	a.get("/jet-page", a.Handlers.JetPage)
	a.get("/sessions", a.Handlers.SessionTest)

	a.get("/users/login", a.Handlers.UserLogin)
	a.post("/users/login", a.Handlers.PostUserLogin)
	a.get("/users/logout", a.Handlers.LogOut)
	a.get("/users/forgot-password", a.Handlers.Forgot)
	a.post("/users/forgot-password", a.Handlers.PostForgot)
	a.get("/users/reset-password", a.Handlers.ResetPasswordForm)
	a.post("/users/reset-password", a.Handlers.PostResetPassword)

	a.get("/form", a.Handlers.Form)
	a.post("/form", a.Handlers.SubmitForm)

	a.get("/json", a.Handlers.JSONExample)
	a.get("/xml", a.Handlers.XMLExample)
	a.get("/download-file", a.Handlers.DownloadExampleFile)

	a.get("/test-crypto", a.Handlers.TestCrypto)

	a.get("/test-cache", a.Handlers.ShowCachePage)
	a.post("/api/save-in-cache", a.Handlers.SaveInCache)
	a.post("/api/get-from-cache", a.Handlers.GetFromCache)
	a.post("/api/delete-from-cache", a.Handlers.DeleteFromCache)
	a.post("/api/empty-cache", a.Handlers.EmptyCache)

	a.get("/test-mail", func(w http.ResponseWriter, r *http.Request) {
		msg := mailer.Message{
			From:        "test@example.com",
			To:          "you@there.com",
			Subject:     "Test Subject - Sent using channel",
			Template:    "test",
			Attachments: nil,
			Data:        nil,
		}

		// //Using Channels
		a.App.Mail.Jobs <- msg
		res := <-a.App.Mail.Results
		if res.Error != nil {
			a.App.ErrorLog.Println(res.Error)
		}

		//You can also call Function Directly
		// err := a.App.Mail.SendMessage_SMTP(msg)
		// if err != nil {
		// 	a.App.ErrorLog.Println(err)
		// }

		fmt.Fprint(w, "Message Sent!")
	})

	a.App.Routes.Get("/create-user", func(w http.ResponseWriter, r *http.Request) {
		u := data.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@here.com",
			Active:    1,
			Password:  "password",
		}

		id, err := a.Models.Users.Insert(u)
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}

		fmt.Fprintf(w, "%d: %s", id, u.FirstName)
	})

	a.App.Routes.Get("/get-all-users", func(w http.ResponseWriter, r *http.Request) {
		users, err := a.Models.Users.GetAll()
		if err != nil {
			a.App.ErrorLog.Println(err)
			return
		}
		for _, x := range users {
			fmt.Fprintf(w, "%d: %s %s\n", x.ID, x.FirstName, x.LastName)
		}
	})

	a.App.Routes.Get("/get-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		u, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Println(err)
		}

		fmt.Fprintf(w, "%d: %s %s --> %s", u.ID, u.FirstName, u.LastName, u.Email)
	})

	a.App.Routes.Get("/update-user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		u, err := a.Models.Users.Get(id)
		if err != nil {
			a.App.ErrorLog.Println(err)
		}

		u.LastName = ""

		validator := a.App.Validator(nil)
		u.Validate(validator)

		if !validator.Valid() {
			for k, v := range validator.Errors {
				fmt.Fprintf(w, "%s: %s\n", k, v)
			}
			return
		}

		err = u.Update(*u)
		if err != nil {
			a.App.ErrorLog.Println(err)
		}

		fmt.Fprintf(w, "%d: %s %s --> %s", u.ID, u.FirstName, u.LastName, u.Email)
	})

	// Static Routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
