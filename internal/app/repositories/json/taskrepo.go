package json

import (
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"task/internal/app/models"
	"task/internal/app/services/configmanager"
)

type TaskRepo struct {
	config *configmanager.Config
}

func NewTaskRepo(config *configmanager.Config) *TaskRepo {
	return &TaskRepo{
		config: config,
	}
}

func (r *TaskRepo) getFileName(id string) string {
	return r.config.JsonPathTask + "/" + id + ".json"
}

func (r *TaskRepo) getID() (string, error) {
	fileName := r.config.JsonPathTask + "/maxid"
	fl, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	maxId, err := strconv.Atoi(string(fl[:]))
	if err != nil {
		return "", err
	}
	maxId++
	err = os.WriteFile(fileName, []byte(strconv.Itoa(maxId)), 0666)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(maxId), nil
}

func (r *TaskRepo) getTaskFromFile(fileName string) (*models.Task, error) {
	var task models.Task
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepo) Save(task *models.Task) (string, error) {
	task.BeforeSave()
	var err error
	if task.ID == "" {
		task.ID, err = r.getID()
		if err != nil {
			return "", err
		}
	}
	bytes, err := json.Marshal(task)
	err = os.WriteFile(r.getFileName(task.ID), bytes, 0666)
	if err != nil {
		return "", err
	}
	return task.ID, nil
}

func (r *TaskRepo) Delete(taskID string) error {
	err := os.Remove(r.getFileName(taskID))
	return err
}

func (r *TaskRepo) GetByID(id string) (*models.Task, error) {
	return r.getTaskFromFile(r.getFileName(id))
}

func (r *TaskRepo) GetAll() ([]models.Task, error) {
	files, err := os.ReadDir(r.config.JsonPathTask)
	if err != nil {
		return nil, err
	}
	var result []models.Task
	for _, file := range files {
		fileName := file.Name()
		ok, _ := regexp.Match("^\\d+\\.json$", []byte(fileName))
		if ok {
			task, err := r.getTaskFromFile(r.config.JsonPathTask + "/" + fileName)
			if err != nil {
				return nil, err
			}
			result = append(result, *task)
		}
	}

	return result, nil
}
