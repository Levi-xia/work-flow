package service

import (
	"encoding/json"
	"workflow/internal/form/model"
	"workflow/internal/form/store"
)

type FormInstance struct {
	Meta  *model.FormInstanceModel
	store store.FormInstanceStore
}

// 创建表单实例
func NewFormInstance(formDefineID int, formData map[string]interface{}) (*FormInstance, error) {
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

// 写入表单数据
func UpdateFormInstanceFormData(formInstanceID int, formData map[string]interface{}) error {
	formDataBytes, err := json.Marshal(formData)
	if err != nil {
		return err
	}
	return store.GetFormInstanceStore().UpdateFormInstanceFormData(formInstanceID, string(formDataBytes))
}
