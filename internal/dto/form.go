package dto

import "workflow/internal/common"

type CreateFormDefineRequest struct {
	Code               string `json:"code" form:"code" binding:"required"`
	Name               string `json:"name" form:"name" binding:"required"`
	FormStructure      string `json:"form_structure" form:"form_structure" binding:"required"`
	ComponentStructure string `json:"component_structure" form:"component_structure" binding:"required"`
}

func (CreateFormDefineRequest) GetMessages() common.ValidatorMessages {
	return common.ValidatorMessages{
		"code.required":                "编码不能为空",
		"name.required":                "名称不能为空",
		"form_structure.required":      "表单结构不能为空",
		"component_structure.required": "组件结构不能为空",
	}
}

type CreateFormDefineResponse struct {
	FormDefineId int `json:"form_define_id"`
}
