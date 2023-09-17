package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	ctnr "github.com/docker/docker/api/types/container"
	cli "github.com/docker/docker/client"
)

type Orchestrator struct {
	Context context.Context
  Cache Cache
	Jobs    []Job
}

func NewOrchestrator(o Orchestrator) (Orchestrator, error) {
	client, err := cli.NewClientWithOpts(cli.FromEnv)

	if err != nil {
		return Orchestrator{}, err
	}

  cache := NewCache(Cache{})

  o.Cache = cache

	ticker := time.NewTicker(5 * time.Second)

	err = o.reattach(*client)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			<-ticker.C
			o.autoHeal(*client)
			o.autoClean(*client)
		}
	}()

	return o, nil
}

func (o *Orchestrator) reattach(client cli.Client) error {
  containers, err := client.ContainerList(o.Context, types.ContainerListOptions{}) 

  if err != nil {
    return err 
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

	cachedJobs, err := o.Cache.ListCache()

	if err != nil {
		return err
	}

	for _, cacheJob := range cachedJobs {
		active := false

		for _, container := range containers {
			if container.ID == cacheJob.ContainerId {
				active = true

				fmt.Println(active)

				break
			}
		}

		if active { continue }

		fmt.Println(cacheJob)

		job := NewJob(Job{
			Name: cacheJob.Name,
			ImageName: cacheJob.ImageName,	
			Port: cacheJob.Port,
		})

		err = job.Run()

		if err != nil {
			return err
		}

		o.Jobs = append(o.Jobs, job)
	}

  return nil 
}

func (o *Orchestrator) autoHeal(client cli.Client) {
	for i, job := range o.Jobs {
		if len(job.Id) == 0 {
			continue
		}

		container, err := client.ContainerInspect(o.Context, job.Id)

		if err != nil {
			fmt.Println(err)

			if cli.IsErrNotFound(err) {
				o.Jobs = append(o.Jobs[:i], o.Jobs[:i+1]...)
			}

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

func (o *Orchestrator) autoClean(client cli.Client) {
	for i, job := range o.Jobs {
		if len(job.Id) == 0 {
			continue
		}

		container, err := client.ContainerInspect(o.Context, job.Id) 

		if err != nil {
			fmt.Println(err)

			if cli.IsErrNotFound(err) {
				fmt.Println("not found")

				o.Jobs = append(o.Jobs[:i], o.Jobs[:i+1]...)
			}

			continue
		}

		if container.RestartCount >= 6 {
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
}

func (o *Orchestrator) List() []CacheJob {
	var jobs []CacheJob

	for i := range o.Jobs {
		jobs = append(jobs, CacheJob{
			Name: o.Jobs[i].Name,
			ContainerId: o.Jobs[i].Id,
		})
	}

	return jobs 
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
