package main

import (
	"fmt"
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
	fmt.Println("Server is running at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
