package store

import "workflow/internal/form/model"

type FormInstanceStore interface {
	CreateFormInstance(formInstance *model.FormInstanceModel) (int, error)
	GetFormInstance(id int) (*model.FormInstanceModel, error)
	UpdateFormInstanceFormData(formInstanceID int, formData map[string]interface{}) error
}
