package render

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
	JetViews   *jet.Set
	Session    *scs.SessionManager
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
	Error           string
	Flash           string //Goes to the session only for the next request.
}

func (c *Render) defaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Secure = c.Secure
	td.ServerName = c.ServerName
	td.CSRFToken = nosurf.Token(r)
	td.Port = c.Port
	if c.Session.Exists(r.Context(), "userID") {
		td.IsAuthenticated = true
	}
	td.Error = c.Session.PopString(r.Context(), "error") //PopString returns string value from session data with given key and deletes it.
	td.Flash = c.Session.PopString(r.Context(), "flash")
	return td
}

func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {

	switch strings.ToLower(c.Renderer) {
	case "go":
		return c.GoPage(w, r, view, data)
	case "jet":
		return c.JetPage(w, r, view, variables, data)
	default:
	}
	return errors.New("RENDERER in .env is not set to 'go' or 'jet'")
}

// GoPage Renders a standard Go Template
func (c *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", c.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData) //This casts the interface{} to *TemplateData
	}

	err = tmpl.Execute(w, &td)
	if err != nil {
		return err
	}

	return nil
}

// JetPage renders a template using the Jet Template Engine.
func (c *Render) JetPage(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	td = c.defaultData(td, r)

	t, err := c.JetViews.GetTemplate(fmt.Sprintf("%s.jet", view))
	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
