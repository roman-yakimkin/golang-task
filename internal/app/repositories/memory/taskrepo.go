package memory

import (
	"sort"
	"task/internal/app/errors"
	"task/internal/app/models"
)

type TaskRepo struct {
	tasks  map[int]models.Task
	lastID int
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		tasks: make(map[int]models.Task),
	}
}

func (r *TaskRepo) Save(task *models.Task) (int, error) {
	task.BeforeSave()
	if task.ID == 0 {
		r.lastID++
		task.ID = r.lastID
	}
	r.tasks[task.ID] = *task
	return task.ID, nil
}

func (r *TaskRepo) Delete(taskID int) error {
	_, ok := r.tasks[taskID]
	if !ok {
		return errors.ErrTaskNotFound
	}
	delete(r.tasks, taskID)
	return nil
}

func (r *TaskRepo) GetByID(id int) (*models.Task, error) {
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
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result, nil
}
