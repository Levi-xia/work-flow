package store

import "workflow/internal/form/model"

type FormDefineStore interface {
	CreateFormDefine(meta *model.FormDefineModel) (int, error)
	GetFormDefine(id int) (*model.FormDefineModel, error)
	GetFormDefinesByCode(code string) ([]*model.FormDefineModel, error)
}
