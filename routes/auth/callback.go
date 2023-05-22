package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ferretcode-freelancing/hawthorn/routes"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type GithubResponse struct {
	AccessToken string `json:"access_token"`
}

type GithubUser struct {
	Id int `json:"id"`
	Login string `json:"login"`
}

func Callback(w http.ResponseWriter, r *http.Request) error {
	code := r.URL.Query().Get("code")	

	token, err := getCode(code)

	if err != nil {
		http.Error(w, "there was an error logging in with github", http.StatusInternalServerError)

		return err
	}

	session, _ := store.Get(r, "hawthorn")

	session.Values["token"] = token

	user, err := getUser(token)

	if err != nil {
		http.Error(w, "there was an error logging in with github", http.StatusInternalServerError)

		return err
	}

	session.Values["owner"] = user.Id
	session.Values["login"] = user.Login

	err = session.Save(r, w)

	if err != nil {
		http.Error(w, "there was an error logging in with github", http.StatusInternalServerError)

		return err
	}

	w.WriteHeader(200)
	w.Write([]byte("you have been authenticated successfully"))

	return nil
}

func getUser(token string) (GithubUser, error) {
	client := http.Client{}

	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)

	if err != nil {
		return GithubUser{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	res, err := client.Do(req)

	if err != nil {
		return GithubUser{}, err
	}

	var user GithubUser

	parseErr := routes.ProcessBody(res.Body, &user)

	if parseErr != nil {
		return GithubUser{}, parseErr
	}

	return user, nil
}


func getCode(code string) (string, error) {
	client := http.Client{}

	clientId := os.Getenv("GH_CLIENT_ID")

	clientSecret := os.Getenv("GH_CLIENT_SECRET")

	url := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		clientId,
		clientSecret,
		code,
	)

	req, err := http.NewRequest(
		"POST",
		url,
		nil,
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	var githubResponse GithubResponse

	parseErr := routes.ProcessBody(res.Body, &githubResponse)

	if parseErr != nil {
		return "", parseErr
	}

	return githubResponse.AccessToken, nil
}
