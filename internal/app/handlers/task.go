package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"task/internal/app/errors"
	"task/internal/app/interfaces"
	"task/internal/app/models"
)

type TaskController struct {
	store interfaces.Store
	vm    interfaces.VoteManager
}

func NewTaskController(store interfaces.Store, vm interfaces.VoteManager) *TaskController {
	return &TaskController{
		store: store,
		vm:    vm,
	}
}

func (c *TaskController) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	var impTask models.ImportTask
	err = json.Unmarshal(body, &impTask)
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}

	var task models.Task
	task.Import(impTask)
	taskId, err := c.store.Task().Save(&task)
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	returnSuccessResponse(w, r, "task has been updated", struct {
		Id string `json:"id"`
	}{Id: taskId})

	_, err = c.vm.DoRouting(&task)
	if err != nil {
	}
}

func (c *TaskController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, ok := vars["id"]
	if returnErrorResponse(!ok, w, r, http.StatusNotFound, errors.ErrTaskNotFound, "") {
		return
	}
	task, err := c.store.Task().GetByID(taskId)
	if returnErrorResponse(err != nil, w, r, http.StatusNotFound, err, "") {
		return
	}
	returnSuccessResponse(w, r, "task has been updated", task)
}

func (c *TaskController) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := c.store.Task().GetAll()
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	returnSuccessResponse(w, r, "", tasks)
}

func (c *TaskController) Update(w http.ResponseWriter, r *http.Request) {
	profile := r.Context().Value("profile").(MiddlewareProfile)
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	var updTask models.UpdateTask
	err = json.Unmarshal(body, &updTask)
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	if returnErrorResponse(profile.UserID != updTask.ID, w, r, http.StatusForbidden, nil, "No enough right to update task") {
		return
	}
	task, err := c.store.Task().GetByID(updTask.ID)
	if returnErrorResponse(err != nil, w, r, http.StatusNotFound, err, "") {
		return
	}
	task.Title = updTask.Title
	task.Body = updTask.Body
	task.CleanReactions()
	_, err = c.store.Task().Save(task)
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	returnSuccessResponse(w, r, "task has been updated", struct {
		Id string `json:"id"`
	}{Id: task.ID})

	c.vm.DoRouting(task)
}

func (c *TaskController) Delete(w http.ResponseWriter, r *http.Request) {
	profile := r.Context().Value("profile").(MiddlewareProfile)
	vars := mux.Vars(r)
	taskId, ok := vars["id"]
	if returnErrorResponse(!ok, w, r, http.StatusNotFound, errors.ErrTaskNotFound, "") {
		return
	}
	task, err := c.store.Task().GetByID(taskId)
	if returnErrorResponse(err != nil, w, r, http.StatusNoContent, err, "") {
		return
	}
	if returnErrorResponse(profile.UserID != task.AuthorID, w, r, http.StatusForbidden, nil, "No enough right to delete task") {
		return
	}
	err = c.store.Task().Delete(taskId)
	if returnErrorResponse(err != nil, w, r, http.StatusNoContent, err, "") {
		return
	}
	returnSuccessResponse(w, r, "task has been updated", struct {
		Id string `json:"id"`
	}{Id: taskId})
}

func (c *TaskController) Vote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	task, err := c.vm.PerformLink(vars["vote_link"])
	if returnErrorResponse(err != nil, w, r, http.StatusNoContent, err, "") {
		return
	}
	_, err = c.vm.DoRouting(task)
	if returnErrorResponse(err != nil, w, r, http.StatusNoContent, err, "") {
		return
	}
	returnSuccessResponse(w, r, "", nil)
}
