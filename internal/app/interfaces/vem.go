package interfaces

import "task/internal/app/models"

type VoteEventManager interface {
	EventApproved(task *models.Task)
	EventRejected(task *models.Task)
}
