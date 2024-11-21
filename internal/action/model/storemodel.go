package model

import (
	"encoding/json"
	"time"
	"workflow/internal/action/bo"
	"workflow/internal/utils"
)

type ActionDefineModel struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	Version      string `json:"version"`
	Protocol     string `json:"protocol"`
	Content      string `json:"content"`
	InputStructs string `json:"input_structs"`
	OutputChecks string `json:"output_checks"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
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
		Protocol:     this.Protocol,
		Content:      content,
		InputStructs: inputStructs,
		OutputChecks: outputChecks,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}
