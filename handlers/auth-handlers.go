package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"myapp/data"
	"net/http"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/lordwestcott/celeritas/mailer"
	"github.com/lordwestcott/celeritas/urlsigner"
)

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "login", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := h.Models.Users.GetByEmail(email)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	matches, err := user.PasswordMatches(password)
	if err != nil {
		w.Write([]byte("Error validatiiing password"))
		return
	}

	if !matches {
		w.Write([]byte("Invalid password"))
		return
	}

	//did the user check remember me?
	if r.Form.Get("remember") == "remember" {
		randomString := h.randomString(12)
		hasher := sha256.New()
		_, err := hasher.Write([]byte(randomString))
		if err != nil {
			h.App.ErrorStatus(w, http.StatusBadRequest)
		}

		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		rm := data.RememberToken{}
		err = rm.InsertToken(user.ID, sha)
		if err != nil {
			h.App.ErrorStatus(w, http.StatusBadRequest)
		}

		//set a cookie
		expire := time.Now().Add(14 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:     fmt.Sprintf("_%s_remember", h.App.AppName),
			Value:    fmt.Sprintf("%d|%s", user.ID, sha),
			Path:     "/",
			Expires:  expire,
			HttpOnly: true,
			Domain:   h.App.Session.Cookie.Domain,
			Secure:   h.App.Session.Cookie.Secure,
			MaxAge:   int(expire.Unix() - time.Now().Unix()),
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)
		h.App.Session.Put(r.Context(), "remember_token", sha)
	}

	h.App.Session.Put(r.Context(), "userID", user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (h *Handlers) LogOut(w http.ResponseWriter, r *http.Request) {
	if h.App.Session.Exists(r.Context(), "remember_token") {
		rt := data.RememberToken{}
		_ = rt.Delete(h.App.Session.GetString(r.Context(), "remember_token"))
	}

	//Delete cookie
	cookie := http.Cookie{
		Name:     fmt.Sprintf("_%s_remember", h.App.AppName),
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-100 * time.Hour),
		HttpOnly: true,
		Domain:   h.App.Session.Cookie.Domain,
		Secure:   h.App.Session.Cookie.Secure,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

	h.App.Session.RenewToken(r.Context())
	h.App.Session.Remove(r.Context(), "userID")
	h.App.Session.Remove(r.Context(), "remember_token")
	h.App.Session.Destroy(r.Context())
	h.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handlers) Forgot(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "forgot", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("Error Rendering: ", err)
		h.App.Error500(w, r)
	}
}

func (h *Handlers) PostForgot(w http.ResponseWriter, r *http.Request) {
	//parse form
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorStatus(w, http.StatusBadRequest)
	}

	//verify email exists
	var u *data.User
	email := r.Form.Get("email")
	u, err = h.Models.Users.GetByEmail(email)
	if err != nil {
		h.App.ErrorStatus(w, http.StatusBadRequest)
	}

	//create a link to password reset form
	link := fmt.Sprintf("%s/users/reset-password?email=%s", h.App.Server.URL, email)
	sign := urlsigner.Signer{
		Secret: []byte(h.App.EncryptionKey),
	}

	//sign the link
	signedLink := sign.GenerateTokenFromString(link)

	//create the email
	var data struct {
		Link string
	}
	data.Link = signedLink
	msg := mailer.Message{
		To:       u.Email,
		Subject:  "Password Reset",
		Template: "password-reset",
		Data:     data,
		From:     "admin@example.com",
	}

	//send the email
	h.App.Mail.Jobs <- msg
	res := <-h.App.Mail.Results
	if res.Error != nil {
		h.App.ErrorStatus(w, http.StatusBadRequest)
	}

	//redirect user
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (h *Handlers) ResetPasswordForm(w http.ResponseWriter, r *http.Request) {
	//get form values
	email := r.URL.Query().Get("email")
	theURL := r.RequestURI
	testUrl := fmt.Sprintf("%s%s", h.App.Server.URL, theURL)

	//validate the url
	signer := urlsigner.Signer{
		Secret: []byte(h.App.EncryptionKey),
	}

	valid := signer.VerifyToken(testUrl)
	if !valid {
		h.App.ErrorLog.Println("Invalid url")
		h.App.ErrorUnauthorized(w, r)
		return
	}

	//make sure its not expired
	expired := signer.Expired(testUrl, 60)
	if expired {
		h.App.ErrorLog.Println("Link expired")
		h.App.ErrorUnauthorized(w, r)
		return
	}

	//display form
	/*

		NOTE:
		You don't want to pass the email in the form, as then the user could change anyones password.
		Instead, you want to use a hidden field to store the email, in an verifiably encrypted format.

	*/
	encryptedEmail, _ := h.encrypt(email)
	vars := make(jet.VarMap)
	vars.Set("email", encryptedEmail)

	err := h.render(w, r, "reset-password", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("Error Rendering: ", err)
		h.App.Error500(w, r)
	}
}

func (h *Handlers) PostResetPassword(w http.ResponseWriter, r *http.Request) {
	//parse form
	err := r.ParseForm()
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	//get and decrypt email
	email, err := h.decrypt(r.Form.Get("email"))
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	//get user
	var u data.User
	user, err := u.GetByEmail(email)
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	//reset the password
	err = user.ResetPassword(user.ID, r.Form.Get("password"))
	if err != nil {
		h.App.Error500(w, r)
		return
	}

	//redirect user
	h.App.Session.Put(r.Context(), "flash", "Password reset successfully")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}
