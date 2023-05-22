package containers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ferretcode-freelancing/hawthorn/builder"
	"github.com/ferretcode-freelancing/hawthorn/orchestrator"
	"github.com/ferretcode-freelancing/hawthorn/routes"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

type Container struct {
	Name    string `json:"name"`
	RepoURL string `json:"repo_url"`
}

func New(w http.ResponseWriter, r *http.Request, o orchestrator.Orchestrator) error {
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

	err = builder.ExtractRepo(session.Values["owner"].(int), repoName[len(repoName)-1])

	if err != nil {
		http.Error(w, "there was an error extracting your repository", http.StatusInternalServerError)

		return err
	}

	err = builder.Build(session.Values["owner"].(int), repoName[len(repoName)-1])

	if err != nil {
		http.Error(w, "there was an error building your repository", http.StatusInternalServerError)

		return err
	}

	fmt.Println(repoName[len(repoName)-1])

	job := orchestrator.NewJob(orchestrator.Job{
		Name:      repoName[len(repoName)-1],
		ImageName: repoName[len(repoName)-1],
	})

	err = o.New(job)

	if err != nil {
		http.Error(w, "there was an error deploying your repository", http.StatusInternalServerError)

		return err
	}

	w.WriteHeader(200)
	w.Write([]byte("your repository was deployed successfully"))

	return nil
}
