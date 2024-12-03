package bo

import (
	"time"
	"workflow/internal/core/constants"
)

type ProcessDefineBo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	UserID    int       `json:"user_id"`
	Version   int       `json:"version"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProcessInstanceBo struct {
	ID              int                             `json:"id"`
	ProcessDefineID int                             `json:"process_define_id"`
	UserID          int                             `json:"user_id"`
	Status          constants.ProcessInstanceStatus `json:"status"`
	Variables       map[string]interface{}          `json:"variables"`
	CreatedAt       time.Time                       `json:"created_at"`
	UpdatedAt       time.Time                       `json:"updated_at"`
}

type ProcessTaskBo struct {
	ID                int `json:"id"`
	ProcessInstanceID int `json:"process_instance_id"`
	FormInstanceID    int `json:"form_instance_id"`
	Name      string                      `json:"name"`
	Code      string                      `json:"code"`
	UserID    int                         `json:"user_id"`
	Status    constants.ProcessTaskStatus `json:"status"`
	Variables map[string]interface{}      `json:"variables"`
	CreatedAt time.Time                   `json:"created_at"`
	UpdatedAt time.Time                   `json:"updated_at"`
}
