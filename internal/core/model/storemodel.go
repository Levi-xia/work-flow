package model

import (
	"time"
	"workflow/internal/core/constants"
)

type ProcessDefineModel struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" db:"name"`
	Code      string    `json:"code" db:"code"`
	Version   int       `json:"version" db:"version"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ProcessInstanceModel struct {
	ID              int                             `json:"id"`
	ProcessDefineID int                             `json:"process_define_id" db:"process_define_id"`
	Status          constants.ProcessInstanceStatus `json:"status" db:"status"`
	Variables       map[string]interface{}          `json:"variables" db:"variables"`
	CreatedAt       time.Time                       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time                       `json:"updated_at" db:"updated_at"`
}

type ProcessTaskModel struct {
	ID                int                         `json:"id"`
	ProcessInstanceID int                         `json:"process_instance_id" db:"process_instance_id"`
	Name              string                      `json:"name" db:"name"`
	Code              string                      `json:"code" db:"code"`
	Status            constants.ProcessTaskStatus `json:"status" db:"status"`
	Variables         map[string]interface{}      `json:"variables" db:"variables"`
	CreatedAt         time.Time                   `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time                   `json:"updated_at" db:"updated_at"`
}
