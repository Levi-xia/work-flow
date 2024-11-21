package service

import (
	"encoding/json"
	"workflow/internal/core/bo"
	"workflow/internal/core/constants"
	"workflow/internal/core/model"
	"workflow/internal/core/store"
)

type ProcessInstance struct {
	Define *ProcessDefine
	Meta   *bo.ProcessInstanceBo
	store  store.ProcessInstanceStore
}

// 创建实例
func NewProcessInstance(define *ProcessDefine, variables map[string]interface{}) (*ProcessInstance, error) {
	var (
		err            error
		variablesBytes []byte
		instanceModel  *model.ProcessInstanceModel
		instanceBo     *bo.ProcessInstanceBo
	)
	if variablesBytes, err = json.Marshal(variables); err != nil {
		return nil, err
	}
	instanceModel = &model.ProcessInstanceModel{
		ProcessDefineID: define.Meta.ID,
		Status:          string(constants.PROCESSINSTANCESTATUSDOING),
		Variables:       string(variablesBytes),
	}
	if instanceBo, err = instanceModel.ToBo(); err != nil {
		return nil, err
	}
	instance := &ProcessInstance{
		Define: define,
		store:  store.GetProcessInstanceStore(),
		Meta:   instanceBo,
	}
	if instance.Meta.ID, err = instance.store.CreateProcessInstance(instanceModel); err != nil {
		return nil, err
	}
	return instance, nil
}

func GetProcessInstance(processInstanceId int) (*bo.ProcessInstanceBo, error) {
	instanceModel, err := store.GetProcessInstanceStore().GetProcessInstance(processInstanceId)
	if err != nil {
		return nil, err
	}
	return instanceModel.ToBo()
}

// 结束实例
func FinishProcessInstance(instanceId int) error {
	return store.GetProcessInstanceStore().FinishProcessInstance(instanceId)
}

func UpdateProcessInstanceVariables(instanceId int, variables map[string]interface{}) error {
	return store.GetProcessInstanceStore().UpdateProcessInstanceVariables(instanceId, variables)
}

// 取消实例
func CancelProcessInstance(instanceId int) error {
	return store.GetProcessInstanceStore().CancelProcessInstance(instanceId)
}
