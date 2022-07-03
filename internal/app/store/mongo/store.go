package mongo

import (
	"task/internal/app/interfaces"
	"task/internal/app/repositories/mongo"
	"task/internal/app/services/dbclient"
)

type Store struct {
	db       *dbclient.MongoDBClient
	taskRepo *mongo.TaskRepo
}

func NewStore(db *dbclient.MongoDBClient, taskRepo *mongo.TaskRepo) *Store {
	return &Store{
		db:       db,
		taskRepo: taskRepo,
	}
}

func (s *Store) Task() interfaces.TaskRepo {
	if s.taskRepo == nil {
		s.taskRepo = mongo.NewTaskRepo(s.db)
	}
	return s.taskRepo
}
