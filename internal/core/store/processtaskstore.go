package store

import "workflow/internal/core/model"

type ProcessTaskStore interface {
	CreateProcessTask(meta *model.ProcessTaskModel) (int, error)
	FinishProcessTask(id int) error
	GetProcessTask(id int) (*model.ProcessTaskModel, error)
	GetRunningTasks(instanceID int) ([]*model.ProcessTaskModel, error)
}
