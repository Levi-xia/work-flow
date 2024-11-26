package dto

import (
	"workflow/internal/action/bo"
	"workflow/internal/common"
)

type CreateActionDefineRequest struct {
	Code         string             `json:"code" form:"code" binding:"required"`
	Name         string             `json:"name" form:"name" binding:"required"`
	Protocol     string             `json:"protocol" form:"protocol" binding:"protocol"`
	Content      *bo.ActionContent  `json:"content" form:"content" binding:"required"`
	InputStructs []*bo.ParamsStruct `json:"input_structs" form:"input_structs"`
	OutputChecks []*bo.OutputCheck  `json:"output_checks" form:"output_checks"`
}

func (CreateActionDefineRequest) GetMessages() common.ValidatorMessages {
	return common.ValidatorMessages{
		"code.required":     "编码不能为空",
		"name.required":     "名称不能为空",
		"protocol.protocol": "协议类型不正确",
		"content.required":  "动作内容不能为空",
	}
}

type CreateActionDefineResponse struct {
	ActionDefineID int `json:"action_define_id"`
}
