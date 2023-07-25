package celeritas

import (
	"net/http"
	"strconv"

	"github.com/justinas/nosurf"
)

// Loading and saving session on each request.
func (c *Celeritas) SessionLoad(next http.Handler) http.Handler {
	c.InfoLog.Println("SessionLoad Called!")
	return c.Session.LoadAndSave(next)
}

func (c *Celeritas) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(c.config.cookie.secure)

	// This will exempt the /someapi/* route from CSRF checks.
	// csrfHandler.ExemptGlob("/someapi/*")
	csrfHandler.ExemptGlob("/api/*")
	//We need to set this as exempt as when we are using the API we are not using the templates

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   c.config.cookie.domain,
	})

	return csrfHandler
}
