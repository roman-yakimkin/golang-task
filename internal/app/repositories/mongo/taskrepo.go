package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"task/internal/app/errors"
	"task/internal/app/models"
	"task/internal/app/services/dbclient"
)

type TaskRepo struct {
	db *dbclient.MongoDBClient
}

func NewTaskRepo(db *dbclient.MongoDBClient) *TaskRepo {
	return &TaskRepo{
		db: db,
	}
}

func (r *TaskRepo) Save(task *models.Task) (string, error) {
	ctx := context.Background()
	client, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return "", err
	}

	task.BeforeSave()
	var taskMongo Task
	err = taskMongo.Import(task)
	if err != nil {
		return "", err
	}
	c := client.Database("task_service").Collection("tasks")
	filter := bson.D{{"_id", taskMongo.ID}}
	insData := bson.M{
		"authorId":  taskMongo.AuthorID,
		"title":     taskMongo.Title,
		"body":      taskMongo.Body,
		"created":   taskMongo.Created,
		"edited":    taskMongo.Edited,
		"reactions": taskMongo.Reactions,
	}
	update := bson.D{{"$set", insData}}

	found := c.FindOne(ctx, filter)
	if found.Err() == mongo.ErrNoDocuments {
		res, err := c.InsertOne(ctx, insData)
		if err != nil {
			return "", err
		}
		task.ID = res.InsertedID.(primitive.ObjectID).Hex()
	} else {
		_, err := c.UpdateOne(ctx, filter, update)
		if err != nil {
			return "", err
		}
	}
	return task.ID, nil
}

func (r *TaskRepo) Delete(taskID string) error {
	ctx := context.Background()
	client, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return err
	}

	c := client.Database("task_service").Collection("tasks")
	id, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}

	one, err := c.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if one.DeletedCount == 0 {
		return errors.ErrOnTaskDeleting
	}
	return nil
}

func (r *TaskRepo) GetByID(taskID string) (*models.Task, error) {
	ctx := context.Background()
	client, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return nil, err
	}

	id, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, err
	}

	c := client.Database("task_service").Collection("tasks")
	result := c.FindOne(ctx, bson.M{"_id": id})

	var mongoTask Task
	err = result.Decode(&mongoTask)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = errors.ErrTaskNotFound
		}
		return nil, err
	}

	task := mongoTask.Export()
	return task, nil
}

func (r *TaskRepo) GetAll() ([]models.Task, error) {
	ctx := context.Background()
	client, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return nil, err
	}

	c := client.Database("task_service").Collection("tasks")
	cursor, err := c.Find(ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	var mongoTasks []Task
	if err = cursor.All(ctx, &mongoTasks); err != nil {
		return nil, err
	}
	results := make([]models.Task, 0, len(mongoTasks))
	for _, mt := range mongoTasks {
		results = append(results, *mt.Export())
	}

	return results, nil
}
