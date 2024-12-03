package model

import (
	"encoding/json"
	"time"
	"workflow/internal/core/bo"
	"workflow/internal/core/constants"
	"workflow/internal/utils"
)

type ProcessDefineModel struct {
	ID        int    `json:"id"`
	Name      string `json:"name" db:"name"`
	Code      string `json:"code" db:"code"`
	UserID    int    `json:"user_id" db:"user_id"`
	Version   int    `json:"version" db:"version"`
	Content   string `json:"content" db:"content"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

func (this *ProcessDefineModel) ToBo() (*bo.ProcessDefineBo, error) {
	var (
		createdAt, updatedAt time.Time
		err                  error
	)
	if createdAt, err = utils.ParseTime(this.CreatedAt); err != nil {
		return nil, err
	}
	if updatedAt, err = utils.ParseTime(this.UpdatedAt); err != nil {
		return nil, err
	}
	return &bo.ProcessDefineBo{
		ID:        this.ID,
		Name:      this.Name,
		Code:      this.Code,
		UserID:    this.UserID,
		Version:   this.Version,
		Content:   this.Content,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

type ProcessInstanceModel struct {
	ID              int    `json:"id"`
	ProcessDefineID int    `json:"process_define_id" db:"process_define_id"`
	UserID          int    `json:"user_id" db:"user_id"`
	Status          string `json:"status" db:"status"`
	Variables       string `json:"variables" db:"variables"`
	CreatedAt       string `json:"created_at" db:"created_at"`
	UpdatedAt       string `json:"updated_at" db:"updated_at"`
}

func (this *ProcessInstanceModel) ToBo() (*bo.ProcessInstanceBo, error) {
	var (
		createdAt, updatedAt time.Time
		variables            map[string]interface{}
		err                  error
	)
	if createdAt, err = utils.ParseTime(this.CreatedAt); err != nil {
		return nil, err
	}
	if updatedAt, err = utils.ParseTime(this.UpdatedAt); err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(this.Variables), &variables); err != nil {
		return nil, err
	}
	return &bo.ProcessInstanceBo{
		ID:              this.ID,
		ProcessDefineID: this.ProcessDefineID,
		UserID:          this.UserID,
		Status:          constants.ProcessInstanceStatus(this.Status),
		Variables:       variables,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}

type ProcessTaskModel struct {
	ID                int    `json:"id"`
	ProcessInstanceID int    `json:"process_instance_id" db:"process_instance_id"`
	FormInstanceID    int    `json:"form_instance_id" db:"form_instance_id"`
	Name              string `json:"name" db:"name"`
	Code              string `json:"code" db:"code"`
	UserID            int    `json:"user_id" db:"user_id"`
	Status            string `json:"status" db:"status"`
	Variables         string `json:"variables" db:"variables"`
	CreatedAt         string `json:"created_at" db:"created_at"`
	UpdatedAt         string `json:"updated_at" db:"updated_at"`
}

func (this *ProcessTaskModel) ToBo() (*bo.ProcessTaskBo, error) {
	var (
		createdAt, updatedAt time.Time
		variables            map[string]interface{}
		err                  error
	)
	if createdAt, err = utils.ParseTime(this.CreatedAt); err != nil {
		return nil, err
	}
	if updatedAt, err = utils.ParseTime(this.UpdatedAt); err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(this.Variables), &variables); err != nil {
		return nil, err
	}
	return &bo.ProcessTaskBo{
		ID:                this.ID,
		ProcessInstanceID: this.ProcessInstanceID,
		FormInstanceID:    this.FormInstanceID,
		Name:              this.Name,
		Code:              this.Code,
		UserID:            this.UserID,
		Status:            constants.ProcessTaskStatus(this.Status),
		Variables:         variables,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}, nil
}
