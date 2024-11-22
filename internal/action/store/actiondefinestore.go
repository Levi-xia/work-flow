package store

import "workflow/internal/action/model"

type ActionDefineStore interface {
	CreateActionDefine(defineModel *model.ActionDefineModel) (int, error)
	GetActionDefine(id int) (*model.ActionDefineModel, error)
	GetActionDefines(ids []int) ([]*model.ActionDefineModel, error)
	GetActionDefinesByCode(code string) ([]*model.ActionDefineModel, error)
}
