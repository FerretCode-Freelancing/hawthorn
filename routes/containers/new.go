package containers

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/ferretcode-freelancing/hawthorn/builder"
	"github.com/ferretcode-freelancing/hawthorn/routes"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

type Container struct {
	Name string
	RepoURL string
}

func New(w http.ResponseWriter, r *http.Request) error {
	session, _ := store.Get(r, "hawthorn")

	if session == nil {
		http.Error(w, "you are not authenticated", http.StatusForbidden)

		return errors.New("you are not authenticated")
	}

	container := Container{}	

	err := routes.ProcessBody(r.Body, &container)

	if err != nil {
		http.Error(w, "there was an error processing the new container", http.StatusInternalServerError)

		return err
	}

	err = builder.Download(r, container.RepoURL)

	if err != nil {
		http.Error(w, "there was an error downloading your repository", http.StatusInternalServerError)

		return err
	}

	repoName := strings.Split(container.RepoURL, "/")

	err = builder.ExtractRepo(session.Values["owner"].(int64), repoName[len(repoName)-1])

	if err != nil {
		http.Error(w, "there was an error extracting your repository", http.StatusInternalServerError)

		return err
	}

	err = builder.Build(session.Values["owner"].(int64), repoName[len(repoName)-1])

	if err != nil {
		http.Error(w, "there was an error building your repository", http.StatusInternalServerError)

		return err
	}

	return nil
}