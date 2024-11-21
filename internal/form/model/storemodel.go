package model

import "time"

type FormDefineModel struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name" db:"name"`
	Code               string    `json:"code" db:"code"`
	FormStructure      string    `json:"form_structure" db:"form_structure"`
	ComponentStructure string    `json:"component_structure" db:"component_structure"`
	Version            int       `json:"version" db:"version"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type FormInstanceModel struct {
	ID           int
	FormDefineID int
	FormData     string    `json:"form_data" db:"form_data"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
