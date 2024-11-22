package model

import (
	"encoding/json"
	"time"
	"workflow/internal/action/bo"
	"workflow/internal/action/constants"
	"workflow/internal/utils"
)

type ActionDefineModel struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Code         string `json:"code" db:"code"`
	Version      int    `json:"version" db:"version"`
	Protocol     string `json:"protocol" db:"protocol"`
	Content      string `json:"content" db:"content"`
	InputStructs string `json:"input_structs" db:"input_structs"`
	OutputChecks string `json:"output_checks" db:"output_checks"`
	CreatedAt    string `json:"created_at" db:"created_at"`
	UpdatedAt    string `json:"updated_at" db:"updated_at"`
}

func (this *ActionDefineModel) ToBo() (*bo.ActionDefineBo, error) {
	var (
		inputStructs []bo.ParamsStruct
		outputChecks []bo.OutputCheck
		content      bo.ActionContent
		createdAt    time.Time
		updatedAt    time.Time
		err          error
	)
	if createdAt, err = utils.ParseTime(this.CreatedAt); err != nil {
		return nil, err
	}
	if updatedAt, err = utils.ParseTime(this.UpdatedAt); err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(this.InputStructs), &inputStructs); err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(this.OutputChecks), &outputChecks); err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(this.Content), &content); err != nil {
		return nil, err
	}
	return &bo.ActionDefineBo{
		ID:           this.ID,
		Name:         this.Name,
		Code:         this.Code,
		Version:      this.Version,
		Protocol:     constants.ActionProtocol(this.Protocol),
		Content:      content,
		InputStructs: inputStructs,
		OutputChecks: outputChecks,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

type ActionRecordModel struct {
	ID             int    `json:"id" db:"id"`
	ActionDefineID int    `json:"action_define_id" db:"action_define_id"`
	ProcessTaskID  int    `json:"process_task_id" db:"process_task_id"`
	Input          string `json:"input" db:"input"`
	Output         string `json:"output" db:"output"`
	CreatedAt      string `json:"created_at" db:"created_at"`
	UpdatedAt      string `json:"updated_at" db:"updated_at"`
}

func (this *ActionRecordModel) ToBo() (*bo.ActionRecordBo, error) {
	var (
		input     map[string]interface{}
		output    map[string]interface{}
		createdAt time.Time
		updatedAt time.Time
		err       error
	)
	if createdAt, err = utils.ParseTime(this.CreatedAt); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(this.Input), &input); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(this.Output), &output); err != nil {
		return nil, err
	}
	return &bo.ActionRecordBo{
		ID:             this.ID,
		ActionDefineID: this.ActionDefineID,
		ProcessTaskID:  this.ProcessTaskID,
		Input:          input,
		Output:         output,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}
