package service

import (
	"encoding/json"
	"workflow/internal/form/bo"
	"workflow/internal/form/model"
	"workflow/internal/form/store"
)

type FormInstance struct {
	Meta  *model.FormInstanceModel
	store store.FormInstanceStore
}

// 创建表单实例
func NewFormInstance(formDefineID int, formData []*bo.FormData) (*FormInstance, error) {
	formDataBytes, err := json.Marshal(formData)
	if err != nil {
		return nil, err
	}
	formInstance := &FormInstance{
		Meta: &model.FormInstanceModel{
			FormDefineID: formDefineID,
			FormData:     string(formDataBytes),
		},
		store: store.GetFormInstanceStore(),
	}
	if formInstance.Meta.ID, err = store.GetFormInstanceStore().CreateFormInstance(formInstance.Meta); err != nil {
		return nil, err
	}
	return formInstance, nil
}

func GetFormInstance(id int) (*bo.FormInstanceBo, error) {
	formInstanceModel, err := store.GetFormInstanceStore().GetFormInstance(id)
	if err != nil {
		return nil, err
	}
	return formInstanceModel.ToBo()
}

// 写入表单数据
func UpdateFormInstanceFormData(formInstanceID int, variables map[string]interface{}) error {
	var (
		err          error
		formData     []*bo.FormData
		formDefine   *bo.FormDefineBo
		formInstance *bo.FormInstanceBo
	)
	// 查询表单定义
	if formInstance, err = GetFormInstance(formInstanceID); err != nil {
		return err
	}
	if formDefine, err = GetFormDefine(formInstance.FormDefineID); err != nil {
		return err
	}
	for _, v := range formDefine.ComponentStructure {
		if value, ok := variables[v.Field]; ok {
			formData = append(formData, &bo.FormData{
				Field: v.Field,
				Value: value,
			})
		}
	}
	formDataBytes, err := json.Marshal(formData)
	if err != nil {
		return err
	}
	return store.GetFormInstanceStore().UpdateFormInstanceFormData(formInstanceID, string(formDataBytes))
}
