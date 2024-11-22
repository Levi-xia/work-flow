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
	FormStructure      string `json:"form_structure" db:"form_structure"`
	ComponentStructure string `json:"component_structure" db:"component_structure"`
	Version            int    `json:"version" db:"version"`
	CreatedAt          string `json:"created_at" db:"created_at"`
	UpdatedAt          string `json:"updated_at" db:"updated_at"`
}

func (this *FormDefineModel) ToBo() *bo.FormDefineBo {
	var (
		formStructure []bo.FormStructure
		createdAt     time.Time
		updatedAt     time.Time
		err           error
	)
	if createdAt, err = utils.ParseTime(this.CreatedAt); err != nil {
		return nil
	}
	if updatedAt, err = utils.ParseTime(this.UpdatedAt); err != nil {
		return nil
	}
	if err = json.Unmarshal([]byte(this.FormStructure), &formStructure); err != nil {
		return nil
	}
	return &bo.FormDefineBo{
		ID:                 this.ID,
		Name:               this.Name,
		Code:               this.Code,
		FormStructure:      formStructure,
		ComponentStructure: this.ComponentStructure,
		Version:            this.Version,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}
}

type FormInstanceModel struct {
	ID           int
	FormDefineID int
	FormData     string `json:"form_data" db:"form_data"`
	CreatedAt    string `json:"created_at" db:"created_at"`
	UpdatedAt    string `json:"updated_at" db:"updated_at"`
}

func (this *FormInstanceModel) ToBo() *bo.FormInstanceBo {
	var (
		formData  []bo.FormData
		createdAt time.Time
		updatedAt time.Time
		err       error
	)
	if createdAt, err = utils.ParseTime(this.CreatedAt); err != nil {
		return nil
	}
	if updatedAt, err = utils.ParseTime(this.UpdatedAt); err != nil {
		return nil
	}
	if err = json.Unmarshal([]byte(this.FormData), &formData); err != nil {
		return nil
	}
	return &bo.FormInstanceBo{
		ID:           this.ID,
		FormDefineID: this.FormDefineID,
		FormData:     formData,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}
