package dto

import "workflow/internal/common"

type CreateProcessDefineRequest struct {
	Content string `json:"content" form:"content" binding:"required"`
}

type CreateProcessDefineResponse struct {
	ProcessDefineId int    `json:"processDefineId"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	Version         int    `json:"version"`
}

func (CreateProcessDefineRequest) GetMessages() common.ValidatorMessages {
	return common.ValidatorMessages{
		"content.required": "内容不能为空",
	}
}

type StartProcessInstanceRequest struct {
	ProcessDefineId int `json:"processDefineId" form:"processDefineId" binding:"required"`
	Variables       map[string]interface{} `json:"variables" form:"variables"`
}

type StartProcessInstanceResponse struct {
	ProcessInstanceId int `json:"processInstanceId"`
}

func (StartProcessInstanceRequest) GetMessages() common.ValidatorMessages {
	return common.ValidatorMessages{
		"processDefineId.required": "流程定义id不能为空",
	}
}

type ExecuteTaskRequest struct {
	TaskId int `json:"taskId" form:"taskId" binding:"required"`
	Variables map[string]interface{} `json:"variables" form:"variables"`
}

type ExecuteTaskResponse struct {
	TaskId int `json:"taskId"`
}

func (ExecuteTaskRequest) GetMessages() common.ValidatorMessages {
	return common.ValidatorMessages{
		"taskId.required": "任务id不能为空",
	}
}
