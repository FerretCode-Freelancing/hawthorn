package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ferretcode-freelancing/hawthorn/routes"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type GithubResponse struct {
	AccessToken string `json:"access_token"`
}

type GithubUser struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

func Callback(w http.ResponseWriter, r *http.Request, deviceCode string, interval int) error {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				err := poll(w, r, deviceCode)

				if err != nil {
					quit <- struct{}{}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func poll(w http.ResponseWriter, r *http.Request, deviceCode string) error {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"https://github.com/login/oauth/access_token?client_id=%s&device_code=%s&grant_type=urn:ietf:params:oauth:grant-type:device_code",
			os.Getenv("GH_CLIENT_ID"),
			deviceCode,
		),
		nil,
	)

	if err != nil {
		return err
	}

	authenticated := false

	githubResponse := GithubResponse{}

	for !authenticated {
		res, err := http.DefaultClient.Do(req)

		if err != nil {
			return err
		}

		tokenResponseBytes, err := io.ReadAll(res.Body)

		if err != nil {
			return err
		}

		if err := json.Unmarshal(tokenResponseBytes, &githubResponse); err != nil {
			return err
		}

		if githubResponse.AccessToken == "" { return nil }

		authenticated = true
	}

	session, _ := store.Get(r, "hawthorn")

	session.Values["token"] = githubResponse.AccessToken

	user, err := getUser(githubResponse.AccessToken)

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