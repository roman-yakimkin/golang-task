package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"task/internal/app/models"
	"time"
)

type Task struct {
	ID        primitive.ObjectID    `bson:"_id, omitempty"`
	AuthorID  string                `bson:"authorId"`
	Title     string                `bson:"title"`
	Body      string                `bson:"body"`
	Created   time.Time             `bson:"created"`
	Edited    time.Time             `bson:"edited"`
	Reactions []models.TaskReaction `bson:"reactions"`
}

func (t *Task) Export() *models.Task {
	var task models.Task
	task.ID = t.ID.Hex()
	task.AuthorID = t.AuthorID
	task.Title = t.Title
	task.Body = t.Body
	task.Created = t.Created
	task.Edited = t.Edited
	task.Reactions = append(task.Reactions, t.Reactions...)
	return &task
}

func (t *Task) Import(task *models.Task) error {
	var err error
	t.ID = primitive.ObjectID{}
	if task.ID != "" {
		t.ID, err = primitive.ObjectIDFromHex(task.ID)
		if err != nil {
			return err
		}
	}
	t.AuthorID = task.AuthorID
	t.Title = task.Title
	t.Body = task.Body
	t.Created = task.Created
	t.Edited = task.Edited
	t.Reactions = append(t.Reactions, task.Reactions...)
	return nil
}
