package controller

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type GoogleController interface {
	IndexHandler(w http.ResponseWriter, r *http.Request)
	GoogleLoginHandler(w http.ResponseWriter, r *http.Request)
	// GoogleCallbackHandler(w http.ResponseWriter, r *http.Request)
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

// func (*googleController) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
// 	profile, err := GetGoogleProfile(r)
// 	if err != nil {
// 		fmt.Fprintf(w, "Failed to get Google profile: %v", err)
// 		return
// 	}
// 	fmt.Fprintf(w, "Hello, %s!", profile.DisplayName)
// }
