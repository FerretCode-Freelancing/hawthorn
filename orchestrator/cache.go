package orchestrator

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Cache struct {
	Path string
}

type CacheData struct {
	Jobs []CacheJob `json:"jobs"`
}

type CacheJob struct {
	Name string `json:"name"`
	ContainerId string `json:"container_id"`
}

func NewCache(c Cache) Cache {
	if c.Path == "" {
		c.Path = "/tmp/hawthorn/cache.json"
	}

	return c
}

func (c *Cache) SearchCache(id string) (Job, error) {
	file, err := os.Open(c.Path)

	if err != nil {
		return Job{}, err
	}

	data, err := io.ReadAll(file)

	if err != nil {
		return Job{}, err
	}

	cacheData := CacheData{}

	err = json.Unmarshal(data, &cacheData)

	if err != nil {
		return Job{}, err
	}

  for _, job := range cacheData.Jobs {
    if job.ContainerId == id {
      job := NewJob(
        Job{
          Name: job.Name,
          ImageName: job.Name,
          Id: job.ContainerId,
          Health: Healthy,
        },
      )

      return job, nil
    }
  }

  return Job{}, errors.New("no job found")
}

func (c *Cache) CacheJob(j Job) error {
	file, err := os.Open(c.Path)

	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	cacheData := CacheData{}

	err = json.Unmarshal(data, &cacheData)

	if err != nil {
		return err
	}

	job := CacheJob{
		Name: j.Name,
		ContainerId: j.Id,
	} 

	cacheData.Jobs = append(cacheData.Jobs, job)

	stringified, err := json.Marshal(cacheData)

	if err != nil {
		return err
	}

	err = os.WriteFile(c.Path, stringified, os.ModeAppend)

	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) UncacheJob(j Job) error {
  file, err := os.Open(c.Path)

	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	cacheData := CacheData{}

	err = json.Unmarshal(data, &cacheData)

	if err != nil {
		return err
	}

  for i := range cacheData.Jobs {
    if cacheData.Jobs[i].ContainerId == j.Id {
      cacheData.Jobs = append(cacheData.Jobs[:i], cacheData.Jobs[i+1:]...)
    }
  }

  stringified, err := json.Marshal(cacheData)

	if err != nil {
		return err
	}

	err = os.WriteFile(c.Path, stringified, os.ModeAppend)

	if err != nil {
		return err
	}

	return nil
}
