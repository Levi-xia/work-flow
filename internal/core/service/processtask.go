package service

import (
	"workflow/internal/core/constants"
	"workflow/internal/core/model"
	"workflow/internal/core/store"
)

type ProcessTask struct {
	Meta     *model.ProcessTaskModel
	Instance *ProcessInstance
	Store    store.ProcessTaskStore
}

// 创建任务
func NewProcessTask(instance *ProcessInstance, code, name string, Variable map[string]interface{}) (*ProcessTask, error) {
	task := &ProcessTask{
		Instance: instance,
		Store:    store.GetProcessTaskStore(),
		Meta: &model.ProcessTaskModel{
			Code:              code,
			Name:              name,
			ProcessInstanceID: instance.Meta.ID,
			Status:            constants.PROCESSTASKSTATUSDOING,
			Variables:         Variable,
		},
	}
	var err error
	if task.Meta.ID, err = task.Store.CreateProcessTask(task.Meta); err != nil {
		return nil, err
	}
	return task, nil
}

func GetRunningTasks(instance *ProcessInstance) ([]*model.ProcessTaskModel, error) {
	store := store.GetProcessTaskStore()
	return store.GetRunningTasks(instance.Meta.ID)
}
