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
