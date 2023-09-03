package orchestrator

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	Healthy   int = 0
	Deploying int = 1
	Error     int = 2
)

type Job struct {
	Context   context.Context
	Name      string
	ImageName string
	Id        string
	Health    int
}

func NewJob(job Job) Job {
	job.Context = context.Background()

	return job
}

func (j *Job) Run() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return err
	}

	fmt.Println(j.ImageName)

	res, err := cli.ContainerCreate(
		j.Context,
		&container.Config{
			Image: j.ImageName,
		},
		nil,
		nil,
		nil,
		"",
	)

	if err != nil {
		return err
	}

	j.Id = res.ID

	if err := cli.ContainerStart(j.Context, res.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}
