package interfaces

import "task/internal/app/models"

type VoteManager interface {
	PerformLink(string) (*models.Task, error)
	DoRouting(*models.Task) (bool, error)
}
