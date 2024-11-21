package bo

import (
	"time"
	"workflow/internal/action/constants"
)

type ActionDefineBo struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	Code         string         `json:"code"`
	Version      string         `json:"version"`
	Protocol     string         `json:"protocol"`
	Content      ActionContent  `json:"content"`
	InputStructs []ParamsStruct `json:"input_structs"`
	OutputChecks []OutputCheck  `json:"output_checks"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type ActionRecordBo struct {
	ID                int                    `json:"id"`
	ActionDefineID    int                    `json:"action_define_id"`
	ProcessInstanceID int                    `json:"process_instance_id"`
	ProcessTaskID     int                    `json:"process_task_id"`
	Input             map[string]interface{} `json:"input"`
	Output            map[string]interface{} `json:"output"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

type ParamsStruct struct {
	Type     constants.ParamsStructType `json:"type"`
	Key      string                     `json:"key"`
	Required bool                       `json:"required"`
}

type OutputCheck struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type ActionContent struct {
	HttpAction HttpAction `json:"http_action"`
}

type HttpAction struct {
	Url     string               `json:"url"`
	Method  constants.HttpMethod `json:"method"`
	Timeout int                  `json:"timeout"`
}
