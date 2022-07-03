package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"task/internal/app/errors"
	"time"
)

const (
	TaskStatusInProcess = iota
	TaskStatusApproved
	TaskStatusRejected
)

type TaskReaction struct {
	Email    string    `json:"email" bson:"email"`
	Sent     time.Time `json:"sent" bson:"sent"`
	Voted    time.Time `json:"voted" bson:"voted"`
	Approved bool      `json:"approved" bson:"approved"`
}

type Task struct {
	ID        string         `json:"id"`
	AuthorID  string         `json:"author_id"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	Created   time.Time      `json:"created"`
	Edited    time.Time      `json:"edited"`
	Reactions []TaskReaction `json:"reactions"`
}

func (t *Task) Import(it ImportTask) {
	t.AuthorID = it.AuthorID
	t.Title = it.Title
	t.Body = it.Body
	for _, email := range it.Emails {
		t.Reactions = append(t.Reactions, TaskReaction{
			Email: email,
		})
	}
}

func (t *Task) BeforeSave() {
	if t.Created.IsZero() {
		t.Created = time.Now()
	}
	t.Edited = time.Now()
}

func (t *Task) Vote(email string, result bool) error {
	idx, reaction, err := t.getReactionByEmail(email)
	if err != nil {
		return err
	}
	if reaction.Sent.IsZero() || !reaction.Voted.IsZero() {
		return errors.ErrVoteLinkNotActive
	}
	t.Reactions[idx].Voted = time.Now()
	t.Reactions[idx].Approved = result
	return nil
}

func (t *Task) CleanReactions() {
	for _, r := range t.Reactions {
		r.Sent = time.Time{}
		r.Voted = time.Time{}
		r.Approved = false
	}
}

func (t *Task) ReactionForSending() (int, *TaskReaction, error) {
	for i, r := range t.Reactions {
		if r.Sent.IsZero() && r.Voted.IsZero() {
			return i, &r, nil
		}
	}
	return 0, nil, errors.ErrTaskReactionNotFound
}

func (t *Task) Status() int {
	for _, r := range t.Reactions {
		if !r.Voted.IsZero() && !r.Approved {
			return TaskStatusRejected
		}
	}
	for _, r := range t.Reactions {
		if !(!r.Sent.IsZero() && !r.Voted.IsZero() && r.Approved) {
			return TaskStatusInProcess
		}
	}
	return TaskStatusApproved
}

func (t *Task) Checksum() string {
	emails := ""
	for _, r := range t.Reactions {
		emails += r.Email
	}
	data := fmt.Sprintf("%s %s %s", t.Title, t.Body, emails)
	result := md5.Sum([]byte(data))
	return hex.EncodeToString(result[:])
}

func (t *Task) getReactionByEmail(email string) (int, *TaskReaction, error) {
	for i, r := range t.Reactions {
		if r.Email == email {
			return i, &r, nil
		}
	}
	return 0, nil, errors.ErrTaskReactionNotFound
}
