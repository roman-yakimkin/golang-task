package errors

import "errors"

var ErrTaskNotFound = errors.New("task not found")
var ErrOnTaskDeleting = errors.New("error on task deleting")

var ErrTaskReactionNotFound = errors.New("task reaction not found")

var ErrVoteLinkNotParsed = errors.New("vote link not parsed")
var ErrVoteLinkNotActive = errors.New("vote link not active")
