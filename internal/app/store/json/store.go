package json

import (
	"task/internal/app/interfaces"
	"task/internal/app/repositories/json"
	"task/internal/app/services/configmanager"
)

type Store struct {
	config   *configmanager.Config
	taskRepo *json.TaskRepo
}

func NewStore(taskRepo *json.TaskRepo, config *configmanager.Config) *Store {
	return &Store{
		config:   config,
		taskRepo: taskRepo,
	}
}

func (s *Store) Task() interfaces.TaskRepo {
	if s.taskRepo == nil {
		s.taskRepo = json.NewTaskRepo(s.config)
	}
	return s.taskRepo
}
