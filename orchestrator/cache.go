package orchestrator

import (
	"encoding/json"
	"errors"
	"fmt"
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
	ImageName string `json:"image_name"`
}

func NewCache(c Cache) Cache {
	if c.Path == "" {
		c.Path = "/tmp/hawthorn/cache.json"
	}

	return c
}

func (c *Cache) ListCache() ([]CacheJob, error) {
	file, err := os.Open(c.Path)

	if err != nil {
		return []CacheJob{}, err
	}

	data, err := io.ReadAll(file)

	if err != nil {
		return []CacheJob{}, err
	}

	cacheData := CacheData{}

	if err := json.Unmarshal(data, &cacheData); err != nil {
		return []CacheJob{}, err
	}

	return cacheData.Jobs, nil
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
		fmt.Println(job)

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
		ImageName: j.ImageName,
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
