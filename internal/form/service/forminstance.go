package service

import (
	"encoding/json"
	"workflow/internal/form/model"
	"workflow/internal/form/store"
)

type FormInstance struct {
	Meta  *model.FormInstanceModel
	Store store.FormInstanceStore
}

func NewFormInstance(formDefineID int, formData map[string]interface{}, store store.FormInstanceStore) (*FormInstance, error) {
	formDataBytes, err := json.Marshal(formData)
	if err != nil {
		return nil, err
	}
	formInstance := &FormInstance{
		Meta: &model.FormInstanceModel{
			FormDefineID: formDefineID,
			FormData:     string(formDataBytes),
		},
	}
	if formInstance.Meta.ID, err = formInstance.Store.CreateFormInstance(formInstance.Meta); err != nil {
		return nil, err
	}
	return formInstance, nil
}
