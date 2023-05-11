package orchestrator

import (
	"context"

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

	if err := cli.ContainerStart(j.Context, res.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(j.Context, res.ID, container.WaitConditionNotRunning)

	// TODO: return channels for orchestrator to check on container status
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	return nil
}
