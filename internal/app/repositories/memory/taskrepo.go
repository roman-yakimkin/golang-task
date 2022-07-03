package memory

import (
	"strconv"
	"task/internal/app/errors"
	"task/internal/app/models"
)

type TaskRepo struct {
	tasks  map[string]models.Task
	lastID int
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		tasks: make(map[string]models.Task),
	}
}

func (r *TaskRepo) Save(task *models.Task) (string, error) {
	task.BeforeSave()
	if task.ID == "" {
		r.lastID++
		task.ID = strconv.Itoa(r.lastID)
	}
	r.tasks[task.ID] = *task
	return task.ID, nil
}

func (r *TaskRepo) Delete(taskID string) error {
	_, ok := r.tasks[taskID]
	if !ok {
		return errors.ErrTaskNotFound
	}
	delete(r.tasks, taskID)
	return nil
}

func (r *TaskRepo) GetByID(id string) (*models.Task, error) {
	task, ok := r.tasks[id]
	if !ok {
		return nil, errors.ErrTaskNotFound
	}
	return &task, nil
}

func (r *TaskRepo) GetAll() ([]models.Task, error) {
	result := make([]models.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		result = append(result, task)
	}
	return result, nil
}
