package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	ctnr "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Orchestrator struct {
	Context context.Context
  Cache Cache
	Jobs    []Job
}

func NewOrchestrator(o Orchestrator) (Orchestrator, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return Orchestrator{}, err
	}

  cache := NewCache(Cache{})

  o.Cache = cache

	ticker := time.NewTicker(5 * time.Second)

	o.reattach(*client)

	go func() {
		for {
			<-ticker.C
			o.autoHeal(*client)
			o.autoClean(*client)
		}
	}()

	return o, nil
}

func (o *Orchestrator) reattach(client client.Client) {
  containers, err := client.ContainerList(o.Context, types.ContainerListOptions{}) 

  if err != nil {
    fmt.Println(err)

    return 
  }

  for _, container := range containers {
    cacheJob, err := o.Cache.SearchCache(container.ID)

    if err != nil {
			if err.Error() != "no job found" {
				fmt.Println(err)
			}

      continue
    }

    o.Jobs = append(o.Jobs, cacheJob)

		fmt.Println(o.Jobs)
  }

  return 
}

func (o *Orchestrator) autoHeal(client client.Client) {
	for _, job := range o.Jobs {
		if len(job.Id) == 0 {
			continue
		}

		container, err := client.ContainerInspect(o.Context, job.Id)

		if err != nil {
			fmt.Println(err)

			continue
		}

		if !container.State.Running || !container.State.Restarting && container.RestartCount < 6 {
			if err := client.ContainerRestart(o.Context, job.Id, ctnr.StopOptions{}); err != nil {
				fmt.Println(err)

				continue
			}
		}
	}
}

func (o *Orchestrator) autoClean(client client.Client) {
	for i, job := range o.Jobs {
		if len(job.Id) == 0 {
			continue
		}

    err := o.Cache.UncacheJob(job)

    if err != nil {
      fmt.Println(err)

      continue
    }
		
		o.Jobs = append(o.Jobs[:i], o.Jobs[i+1:]...)

		err = client.ContainerRemove(o.Context, job.Id, types.ContainerRemoveOptions{})

		if err != nil {
			fmt.Println(err)

			continue
		}
	}
}

func (o *Orchestrator) List() []Job {
	return o.Jobs
}

func (o *Orchestrator) Get(job string) (Job, error) {
	for i := range o.Jobs {
		if o.Jobs[i].Name == job {
			return o.Jobs[i], nil
		}
	}

	return Job{}, errors.New("not found")
}

func (o *Orchestrator) New(job Job) error {
	o.Jobs = append(o.Jobs, job)

	err := job.Run()

	if err != nil {
		return err
	}

  err = o.Cache.CacheJob(job)

  if err != nil {
    return err
  }



	return nil
}
