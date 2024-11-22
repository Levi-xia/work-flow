package store

import "workflow/internal/action/model"

type ActionRecordStore interface {
	CreateActionRecord(record *model.ActionRecordModel) (int, error)
	GetActionRecord(id int) (*model.ActionRecordModel, error)
	GetActionRecordsByProcessTaskID(processTaskID int) ([]*model.ActionRecordModel, error)
	UpdateActionRecordResult(id int, input map[string]interface{}, output map[string]interface{}) error
}
