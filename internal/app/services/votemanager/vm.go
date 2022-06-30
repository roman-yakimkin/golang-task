package votemanager

import (
	"task/internal/app/errors"
	"task/internal/app/interfaces"
	"task/internal/app/models"
	"time"
)

type VoteManager struct {
	store interfaces.Store
	vlm   interfaces.VoteLinkManager
	vem   interfaces.VoteEventManager
	msm   interfaces.MailSenderManager
}

func NewVoteManager(
	store interfaces.Store,
	vlm interfaces.VoteLinkManager,
	vem interfaces.VoteEventManager,
	msm interfaces.MailSenderManager) *VoteManager {
	return &VoteManager{
		store: store,
		vlm:   vlm,
		vem:   vem,
		msm:   msm,
	}
}

func (vm *VoteManager) PerformLink(link string) (*models.Task, error) {
	data, err := vm.vlm.Parse(link)
	if err != nil {
		return nil, err
	}
	task, err := vm.store.Task().GetByID(data.TaskID)
	if err != nil {
		return nil, err
	}
	if task.Checksum() != data.Checksum {
		return nil, errors.ErrVoteLinkNotActive
	}
	err = task.Vote(data.Email, data.Result)
	if err != nil {
		return nil, err
	}
	_, err = vm.store.Task().Save(task)
	if err != nil {
		return nil, err
	}
	if task.Status() == models.TaskStatusApproved {
		vm.vem.EventApproved(task)
	}
	if task.Status() == models.TaskStatusRejected {
		vm.vem.EventRejected(task)
	}
	return task, nil
}

func (vm *VoteManager) DoRouting(task *models.Task) (bool, error) {
	if task.Status() != models.TaskStatusInProcess {
		return false, nil
	}
	idx, reaction, err := task.ReactionForSending()
	if err != nil {
		return false, err
	}
	linkYes := "vote/" + vm.vlm.Generate(interfaces.VoteLinkData{
		Email:    reaction.Email,
		TaskID:   task.ID,
		Result:   true,
		Checksum: task.Checksum(),
	})
	linkNo := "vote/" + vm.vlm.Generate(interfaces.VoteLinkData{
		Email:    reaction.Email,
		TaskID:   task.ID,
		Result:   false,
		Checksum: task.Checksum(),
	})
	err = vm.msm.DoSendingMail(task, linkYes, linkNo)
	if err != nil {
		return false, err
	}
	task.Reactions[idx].Sent = time.Now()
	_, err = vm.store.Task().Save(task)
	if err != nil {
		return false, err
	}
	return true, nil
}
