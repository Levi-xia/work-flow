package dto

import "workflow/internal/common"


type CreateFormDefineRequest struct {
	Code string `json:"code" form:"code" binding:"required"`
	Name string `json:"name" form:"name" binding:"required"`
	FormStructure string `json:"formStructure" form:"formStructure" binding:"required"`
	ComponentStructure string `json:"componentStructure" form:"componentStructure" binding:"required"`
}

func (CreateFormDefineRequest) GetMessages() common.ValidatorMessages {
	return common.ValidatorMessages{
		"code.required": "编码不能为空",
		"name.required": "名称不能为空",
		"formStructure.required": "表单结构不能为空",
		"componentStructure.required": "组件结构不能为空",
	}
}

type CreateFormDefineResponse struct {
	FormDefineId int `json:"formDefineId"`
}
