package service

import (
	"encoding/json"
	"workflow/internal/core/bo"
	"workflow/internal/core/constants"
	"workflow/internal/core/model"
	"workflow/internal/core/store"
)

type ProcessTask struct {
	Meta     *bo.ProcessTaskBo
	Instance *ProcessInstance
	store    store.ProcessTaskStore
}

// 创建任务
func NewProcessTask(instance *ProcessInstance, code, name string, variables map[string]interface{}) (*ProcessTask, error) {
	var (
		err            error
		variablesBytes []byte
		taskModel      *model.ProcessTaskModel
		taskBo         *bo.ProcessTaskBo
	)
	if variablesBytes, err = json.Marshal(variables); err != nil {
		return nil, err
	}
	taskModel = &model.ProcessTaskModel{
		Code:              code,
		Name:              name,
		ProcessInstanceID: instance.Meta.ID,
		Status:            string(constants.PROCESSTASKSTATUSDOING),
		Variables:         string(variablesBytes),
	}
	if taskBo, err = taskModel.ToBo(); err != nil {
		return nil, err
	}
	task := &ProcessTask{
		Instance: instance,
		store:    store.GetProcessTaskStore(),
		Meta:     taskBo,
	}
	if task.Meta.ID, err = task.store.CreateProcessTask(taskModel); err != nil {
		return nil, err
	}
	return task, nil
}

func GetProcessTask(taskId int) (*bo.ProcessTaskBo, error) {
	taskModel, err := store.GetProcessTaskStore().GetProcessTask(taskId)
	if err != nil {
		return nil, err
	}
	return taskModel.ToBo()
}

func FinishProcessTask(taskId int, variables map[string]interface{}) error {
	return store.GetProcessTaskStore().FinishProcessTask(taskId, variables)
}

func GetRunningTasks(instance *ProcessInstance) ([]*bo.ProcessTaskBo, error) {
	var (
		err        error
		taskModels []*model.ProcessTaskModel
		taskBos    []*bo.ProcessTaskBo
	)
	if taskModels, err = store.GetProcessTaskStore().GetRunningTasks(instance.Meta.ID); err != nil {
		return nil, err
	}
	taskBos = make([]*bo.ProcessTaskBo, 0)
	for _, taskModel := range taskModels {
		taskBo, err := taskModel.ToBo()
		if err != nil {
			return nil, err
		}
		taskBos = append(taskBos, taskBo)
	}
	return taskBos, nil
}
