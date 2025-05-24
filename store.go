package main

import (
	"encoding/json"
	"os"
)

type Store struct {
	filename string
}

func (c *Store) SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, data, 0644)
}

func (s *Store) LoadTasks() ([]Task, error) {
	data, err := os.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}
