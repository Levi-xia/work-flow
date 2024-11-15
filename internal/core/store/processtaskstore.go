package store

import "workflow/internal/core/model"

type ProcessTaskStore interface {
	CreateProcessTask(meta *model.ProcessTaskModel) (int, error)
	FinishProcessTask(id int, variables map[string]interface{}) error
	GetProcessTask(id int) (*model.ProcessTaskModel, error)
	GetRunningTasks(instanceID int) ([]*model.ProcessTaskModel, error)
}
