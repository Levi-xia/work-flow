package store

import "workflow/internal/core/model"


type ProcessInstanceStore interface {
	CreateProcessInstance(meta *model.ProcessInstanceModel) (int, error)
	FinishProcessInstance(id int) error
	CancelProcessInstance(id int) error
	GetProcessInstance(id int) (*model.ProcessInstanceModel, error)
	AddVariables(id int, k string, v interface{}) error
	DeleteVariables(id int, variables ...string) error
}
