package store

import "workflow/internal/core/model"


type ProcessDefineStore interface {
	// 获取定义
	GetProcessDefine(id int) (*model.ProcessDefineModel, error)
	// 创建定义
	CreateProcessDefine(meta *model.ProcessDefineModel) (int, error)
	// 获取定义
	GetProcessDefinesByCode(code string) ([]*model.ProcessDefineModel, error)
}
