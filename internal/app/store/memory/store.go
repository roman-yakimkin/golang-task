package memory

import (
	"task/internal/app/interfaces"
	"task/internal/app/repositories/memory"
)

type Store struct {
	taskRepo *memory.TaskRepo
}

func NewStore(taskRepo *memory.TaskRepo) *Store {
	return &Store{
		taskRepo: taskRepo,
	}
}

func (s *Store) Task() interfaces.TaskRepo {
	if s.taskRepo == nil {
		s.taskRepo = memory.NewTaskRepo()
	}
	return s.taskRepo
}
