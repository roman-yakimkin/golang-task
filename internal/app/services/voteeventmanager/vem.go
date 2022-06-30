package voteeventmanager

import (
	"fmt"
	"task/internal/app/models"
)

type VoteEventManager struct {
}

func NewVoteEventManager() *VoteEventManager {
	return &VoteEventManager{}
}

func (m *VoteEventManager) EventApproved(task *models.Task) {
	fmt.Printf("Task %d approved\n", task.ID)
}

func (m *VoteEventManager) EventRejected(task *models.Task) {
	fmt.Printf("Task %d rejected\n", task.ID)
}
