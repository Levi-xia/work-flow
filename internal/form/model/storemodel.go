package model

import (
	"encoding/json"
	"time"
	"workflow/internal/form/bo"
	"workflow/internal/utils"
)

type FormDefineModel struct {
	ID                 int    `json:"id"`
	Name               string `json:"name" db:"name"`
	Code               string `json:"code" db:"code"`
	UserID             int    `json:"user_id" db:"user_id"`
	FormStructure      string `json:"form_structure" db:"form_structure"`
	ComponentStructure string `json:"component_structure" db:"component_structure"`
	Version            int    `json:"version" db:"version"`
	CreatedAt          string `json:"created_at" db:"created_at"`
	UpdatedAt          string `json:"updated_at" db:"updated_at"`
}

func (this *FormDefineModel) ToBo() (*bo.FormDefineBo, error) {
	var (
		componentStructure []bo.Component
		createdAt          time.Time
		updatedAt          time.Time
		err                error
	)
	if createdAt, err = utils.ParseTime(this.CreatedAt); err != nil {
		return nil, err
	}
	if updatedAt, err = utils.ParseTime(this.UpdatedAt); err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(this.ComponentStructure), &componentStructure); err != nil {
		return nil, err
	}
	return &bo.FormDefineBo{
		ID:                 this.ID,
		Name:               this.Name,
		Code:               this.Code,
		UserID:             this.UserID,
		FormStructure:      this.FormStructure,
		ComponentStructure: componentStructure,
		Version:            this.Version,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}, nil
}

type FormInstanceModel struct {
	ID           int
	FormDefineID int    `json:"form_define_id" db:"form_define_id"`
	FormData     string `json:"form_data" db:"form_data"`
	CreatedAt    string `json:"created_at" db:"created_at"`
	UpdatedAt    string `json:"updated_at" db:"updated_at"`
}

func (this *FormInstanceModel) ToBo() (*bo.FormInstanceBo, error) {
	var (
		formData  []bo.FormData
		createdAt time.Time
		updatedAt time.Time
		err       error
	)
	if createdAt, err = utils.ParseTime(this.CreatedAt); err != nil {
		return nil, err
	}
	if updatedAt, err = utils.ParseTime(this.UpdatedAt); err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(this.FormData), &formData); err != nil {
		return nil, err
	}
	return &bo.FormInstanceBo{
		ID:           this.ID,
		FormDefineID: this.FormDefineID,
		FormData:     formData,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}
