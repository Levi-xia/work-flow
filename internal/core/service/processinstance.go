package service

import (
	"workflow/internal/core/constants"
	"workflow/internal/core/model"
	"workflow/internal/core/store"
)


type ProcessInstance struct {
	Define *ProcessDefine
	Meta   *model.ProcessInstanceModel
	Store  store.ProcessInstanceStore
}

// 创建实例
func NewProcessInstance(define *ProcessDefine, variables map[string]interface{}) (*ProcessInstance, error) {
	instance := &ProcessInstance{
		Define: define,
		Store:  store.GetProcessInstanceStore(),
		Meta: &model.ProcessInstanceModel{
			ProcessDefineID: define.Meta.ID,
			Status:          constants.PROCESSINSTANCESTATUSDOING,
			Variables:       variables,
		},
	}
	var err error
	if instance.Meta.ID, err = instance.Store.CreateProcessInstance(instance.Meta); err != nil {
		return nil, err
	}
	return instance, nil
}

// 结束实例
func (this *ProcessInstance) FinishProcessInstance() error {
	return this.Store.FinishProcessInstance(this.Meta.ID)
}

// 取消实例
func (this *ProcessInstance) CancelProcessInstance() error {
	return this.Store.CancelProcessInstance(this.Meta.ID)
}
