package mailsendingmanager

import (
	"fmt"
	"task/internal/app/models"
)

type DummyMailSendingManager struct {
}

func NewDummyMailSendingManager() *DummyMailSendingManager {
	return &DummyMailSendingManager{}
}

func (m *DummyMailSendingManager) DoSendingMail(task *models.Task, linkYes string, linkNo string) error {
	fmt.Printf("ID: %d, Title: %s\n", task.ID, task.Title)
	fmt.Printf("Yes link: %s\n", linkYes)
	fmt.Printf("No link : %s\n", linkNo)
	return nil
}
