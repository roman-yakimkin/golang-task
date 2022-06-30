package interfaces

import "task/internal/app/models"

type MailSenderManager interface {
	DoSendingMail(*models.Task, string, string) error
}
