package orchestrator

import (
	"encoding/json"
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

func (c *Cache) UncacheJob() error {
	return nil
}