package dto

import "workflow/internal/common"

type CreateProcessDefineRequest struct {
	Content string `json:"content" form:"content" binding:"required"`
}

type CreateProcessDefineResponse struct {
	ProcessDefineId int    `json:"process_define_id"`
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
	ProcessDefineId int `json:"process_define_id" form:"process_define_id" binding:"required"`
	Variables       map[string]interface{} `json:"variables" form:"variables"`
}

type StartProcessInstanceResponse struct {
	ProcessInstanceId int `json:"process_instance_id"`
}

func (StartProcessInstanceRequest) GetMessages() common.ValidatorMessages {
	return common.ValidatorMessages{
		"process_define_id.required": "流程定义id不能为空",
	}
}

type ExecuteTaskRequest struct {
	TaskId int `json:"task_id" form:"task_id" binding:"required"`
	Variables map[string]interface{} `json:"variables" form:"variables"`
}

type ExecuteTaskResponse struct {
	TaskId int `json:"task_id"`
}

func (ExecuteTaskRequest) GetMessages() common.ValidatorMessages {
	return common.ValidatorMessages{
		"task_id.required": "任务id不能为空",
	}
}
