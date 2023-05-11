package builder

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func Build(ownerId int64, repoName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/tmp/hawthorn/out/%s-%s/Dockerfile",
		strconv.FormatInt(ownerId, 10),
		repoName,
	)

	file, err := os.Open(path)

	if err != nil {
		return err
	}

	ctx := context.Background()

	res, err := cli.ImageBuild(ctx, file, types.ImageBuildOptions{})

	if err != nil {
		return err
	}

	res.Body.Close()

	return nil
}