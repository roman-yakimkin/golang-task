package errors

import "errors"

var ErrTaskNotFound = errors.New("task not found")
var ErrOnTaskDeleting = errors.New("error on task deleting")

var ErrTaskReactionNotFound = errors.New("task reaction not found")

var ErrVoteLinkNotParsed = errors.New("vote link not parsed")
var ErrVoteLinkNotActive = errors.New("vote link not active")

var ErrNoEnoughRightsToCreateTask = errors.New("no enough rights to create a task")
var ErrNoEnoughRightsToUpdateTask = errors.New("no enough rights to update a task")
var ErrNoEnoughRightsToDeleteTask = errors.New("no enough rights to delete a task")
