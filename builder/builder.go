package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func Build(ownerId int, repoName string, entrypointPath string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	path := fmt.Sprintf(
		"%s/hawthorn/out/%s-%s",
		homeDir,
		strconv.FormatInt(int64(ownerId), 10),
		repoName,
	)

	directories, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	var currentDirectory os.DirEntry 

	currentDirectory = directories[0]
	currentIndex := 0

	currentInfo, err := currentDirectory.Info()

	if err != nil {
		return err
	}

	for i, directory := range directories {
		info, err := directory.Info()

		if err != nil {
			return err
		}

		if info.ModTime().After(currentInfo.ModTime()) {
			currentDirectory = directory	
			currentIndex = i

			newInfo, err := currentDirectory.Info()

			if err != nil {
				return err
			}

			currentInfo = newInfo
		}
	}

	directories = append(directories[:currentIndex], directories[:currentIndex+1]...)

	for _, directory := range directories {
		err := os.RemoveAll(directory.Name())

		if err != nil {
			return err
		}
	}

	directories = []os.DirEntry{
		currentDirectory,
	}

	path += "/" + directories[0].Name()

	if entrypointPath != "" {
		path += entrypointPath
	}

	buildContext, err := archive.TarWithOptions(path, &archive.TarOptions{})

	if err != nil {
		return err
	}

	defer buildContext.Close()

	ctx := context.Background()

	res, err := cli.ImageBuild(ctx, buildContext, types.ImageBuildOptions{
		Context: buildContext,
		Tags: []string{fmt.Sprintf("%s:latest", repoName)},
		Dockerfile: "Dockerfile",
		Remove: true,
	})

	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, res.Body)

	if err != nil {
		return err
	}

	res.Body.Close()

	return nil
}
