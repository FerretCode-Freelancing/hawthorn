package orchestrator

import (
	"context"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	Healthy   int = 0
	Deploying int = 1
	Error     int = 2
)

type Job struct {
	Context   context.Context
	Port int
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

	var port nat.Port

	if j.Port > 0 {
		port, err = nat.NewPort("tcp", strconv.Itoa(j.Port))

		if err != nil {
			return err
		}
	}

	hostConfig := container.HostConfig{
		PortBindings: nat.PortMap{
			port: []nat.PortBinding{
				{
					HostIP: "0.0.0.0",
					HostPort: strconv.Itoa(j.Port),
				},
			},
		},
	}

	var res container.CreateResponse

	if j.Port > 0 {
		res, err = cli.ContainerCreate(
			j.Context,
			&container.Config{
				Image: j.ImageName,
				ExposedPorts: nat.PortSet{
					port: struct{}{},
				},
			},
			&hostConfig,
			nil,
			nil,
			"",
		)

		if err != nil {
			return err
		}
	} else {
		res, err = cli.ContainerCreate(
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
	}

	if err != nil {
		return err
	}

	j.Id = res.ID

	if err := cli.ContainerStart(j.Context, res.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}
