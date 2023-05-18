package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type GoogleController interface {
	IndexHandler(w http.ResponseWriter, r *http.Request)
	GoogleLoginHandler(w http.ResponseWriter, r *http.Request)
	GoogleCallbackHandler(w http.ResponseWriter, r *http.Request)
}

type Profile struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
}

type controller struct {
	oauthConf *oauth2.Config
}

func NewGoogleController(conf *oauth2.Config) GoogleController {
	return &controller{
		oauthConf: conf,
	}
}

func (c *controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/google/login">GoogleでLogin</a></body></html>`
	fmt.Fprintf(w, html)
}

func (c *controller) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := c.oauthConf.AuthCodeURL("random")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (c *controller) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != "random" {
		fmt.Fprintf(w, "Invalid state parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := c.oauthConf.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Fprintf(w, "Code exchange failed with %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken))
	if err != nil {
		fmt.Fprintf(w, "Failed getting user info: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer resp.Body.Close()

	var profile Profile
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		fmt.Fprintf(w, "Failed reading response body: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("%+v\n", profile)

	jsonBytes, err := json.Marshal(profile)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(jsonBytes)
}
