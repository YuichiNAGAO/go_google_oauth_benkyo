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
	state := "random"
	url := c.oauthConf.AuthCodeURL(state)
	fmt.Println(url)
	cookie := &http.Cookie{
		Name: "token",
		// 認証用のトークン。jwt入れることも多いと思います
		Value: state,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (c *controller) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("token")
	fmt.Println(state.Value)
	fmt.Println(r.FormValue("state"))

	if r.FormValue("state") != state.Value {
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
