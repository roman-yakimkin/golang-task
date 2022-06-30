package interfaces

import "task/internal/app/models"

type TaskRepo interface {
	Save(task *models.Task) (int, error)
	Delete(int) error
	GetByID(id int) (*models.Task, error)
	GetAll() ([]models.Task, error)
}
