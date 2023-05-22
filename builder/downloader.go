package builder

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

func Download(r *http.Request, repoURL string) error {
	ctx := context.Background()

	session, _ := store.Get(r, "hawthorn")

	if session.Values["token"] == nil {
		return errors.New("you are not authenticated")
	}

	fmt.Println(session.Values["token"])

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: session.Values["token"].(string)},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repoName := strings.Split(repoURL, "/")

	repo, _, err := client.Repositories.Get(
		ctx, 
		session.Values["login"].(string),
		repoName[len(repoName) - 1],
	)

	if err != nil {
		return err
	}

	httpClient := &http.Client{}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/zipball/%s", *repo.URL, repo.GetMasterBranch()),
		nil,
	) 

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", session.Values["token"].(string)))

	if err != nil {
		return err
	}

	zipball, err := httpClient.Do(req)

	if err != nil {
		return err
	}

	defer zipball.Body.Close()

	zipballBody, err := io.ReadAll(zipball.Body)

	if err != nil {
		return err
	}

	file, err := os.Create(
		fmt.Sprintf(
			"/tmp/hawthorn/%s-%s.zip",
			strconv.FormatInt(*repo.Owner.ID, 10),
			*repo.Name,
		),
	)

	if err != nil {
		return err
	}

	defer file.Close()

	_, downloadErr := io.Copy(file, bytes.NewReader(zipballBody))

	if downloadErr != nil {
	 return err
	}

	fmt.Printf(
		"Downloaded %s-%s.zip\n",
		strconv.FormatInt(*repo.Owner.ID, 10),
		*repo.Name,
	)

	return nil
}

func ExtractRepo(ownerId int, repoName string) error {
	outputDir := fmt.Sprintf(
		"/tmp/hawthorn/out/%s-%s",
		strconv.FormatInt(int64(ownerId), 10),
		repoName,
	)

	zipball, err := zip.OpenReader(
		fmt.Sprintf(
			"/tmp/hawthorn/%s-%s.zip",
			strconv.FormatInt(int64(ownerId), 10),
			repoName,
		),
	)

	if err != nil {
		return err
	}

	defer zipball.Close()

	for _, file := range zipball.File {
		file.Name = filepath.Base(file.Name)

		path := filepath.Join(outputDir, file.Name)

		if !strings.HasPrefix(path, filepath.Clean(outputDir)+string(os.PathSeparator)) {
			return errors.New("invalid file path")
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)

			continue
		}

		if err := os.MkdirAll(filepath.Dir(filepath.Dir(path)), os.ModePerm); err != nil {
			return err
		}

		// create file in output directory
		destFile, err := os.Create(path)

		if err != nil {
			return err
		}

		// open file in archive
		zipballFile, err := file.Open()

		if err != nil {
			return err
		}

		// copy file in archive to the empty destination file
		if _, err := io.Copy(destFile, zipballFile); err != nil {
			return err
		}

		destFile.Close()
		zipballFile.Close()
	}

	return nil
}
