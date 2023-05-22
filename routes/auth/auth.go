package auth

import (
	"fmt"
	"net/http"
	"os"
)

func Login(w http.ResponseWriter, r *http.Request) error {
	clientId := os.Getenv("GH_CLIENT_ID")

	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&scope=repo,read:user&redirect_uri=http://%s:%s/auth/callback",
		clientId,
		os.Getenv("CALLBACK_URL_HOST"),
		os.Getenv("CALLBACK_URL_PORT"),
	)

	http.Redirect(w, r, url, http.StatusFound)

	return nil
}
