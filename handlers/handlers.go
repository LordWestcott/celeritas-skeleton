package handlers

import (
	"fmt"
	"myapp/data"
	"net/http"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/lordwestcott/celeritas"
)

type Handlers struct {
	App    *celeritas.Celeritas
	Models data.Models
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	defer h.App.LoadTime(time.Now())
	err := h.App.Render.Page(w, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering: ", err)
	}
}

func (h *Handlers) GoPage(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.GoPage(w, r, "home", nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering: ", err)
	}
}

func (h *Handlers) JetPage(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.JetPage(w, r, "jet-template", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering: ", err)
	}
}

func (h *Handlers) SessionTest(w http.ResponseWriter, r *http.Request) {
	//Store data in session
	myData := "bar"
	h.App.Session.Put(r.Context(), "foo", myData)

	//Retrieve data from session
	myValue := h.App.Session.GetString(r.Context(), "foo")

	//Add data to Jet template via a Jet VarMap
	vars := make(jet.VarMap)
	vars.Set("foo", myValue)

	//Pass the VarMap to the JetPage method
	err := h.App.Render.JetPage(w, r, "sessions", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering: ", err)
	}
}

func (h *Handlers) JSONExample(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID      int64    `json:"id"`
		Name    string   `json:"name"`
		Hobbies []string `json:"hobbies"`
	}

	payload.ID = 1
	payload.Name = "John Doe"
	payload.Hobbies = []string{"hiking", "biking", "swimming"}

	err := h.App.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		h.App.ErrorLog.Println("error writing json: ", err)
	}
}

func (h *Handlers) XMLExample(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		ID      int64    `xml:"id"`
		Name    string   `xml:"name"`
		Hobbies []string `xml:"hobbies>hobby"`
	} //interestingly this doesn't work if you create the struct as a var instead of a type.

	pay := Payload{}

	pay.ID = 1
	pay.Name = "John Doe"
	pay.Hobbies = []string{"hiking", "biking", "swimming"}

	err := h.App.WriteXML(w, http.StatusOK, pay)
	if err != nil {
		h.App.ErrorLog.Println("error writing xml: ", err)
	}
}

func (h *Handlers) DownloadExampleFile(w http.ResponseWriter, r *http.Request) {
	err := h.App.DownloadFile(w, r, "./public/images", "celeritas.jpg")
	if err != nil {
		h.App.ErrorLog.Println("error downloading file: ", err)
	}
}

func (h *Handlers) TestCrypto(w http.ResponseWriter, r *http.Request) {
	plainText := "Hello, world!"

	fmt.Fprint(w, "Unencrypted: "+plainText+"\n")
	encrypted, err := h.encrypt(plainText)
	if err != nil {
		h.App.ErrorLog.Println("error encrypting: ", err)
		h.App.Error500(w, r)
	}

	fmt.Fprint(w, "Encrypted: "+encrypted+"\n")
	decrypted, err := h.decrypt(encrypted)
	if err != nil {
		h.App.ErrorLog.Println("error decrypting: ", err)
		h.App.Error500(w, r)
	}

	fmt.Fprint(w, "Decrypted: "+decrypted+"\n")
}
