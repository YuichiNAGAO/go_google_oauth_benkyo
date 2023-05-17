package main

import (
	"net/http"

	"github.com/YuichiNAGAO/go_google_oauth_benkyo/config"
	"github.com/YuichiNAGAO/go_google_oauth_benkyo/controller"
)

var (
	oauthConf        = config.SetUpOauth()
	googleController = controller.NewGoogleController(oauthConf)
)

func main() {
	http.HandleFunc("/", googleController.IndexHandler)
	// http.HandleFunc("/google/login", googleController.GoogleLoginHandler)
	// http.HandleFunc("/google/callback", googleController.GoogleCallbackHandler)
	http.ListenAndServe(":8080", nil)
}
