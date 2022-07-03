package interfaces

import "task/internal/app/models"

type TaskRepo interface {
	Save(task *models.Task) (string, error)
	Delete(string) error
	GetByID(id string) (*models.Task, error)
	GetAll() ([]models.Task, error)
}
