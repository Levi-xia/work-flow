package bo

import "time"

type FormStructure struct {
	Type  string `json:"type"`
	Field string `json:"field"`
	Title string `json:"title"`
}

type FormData struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type FormDefineBo struct {
	ID                 int             `json:"id"`
	Name               string          `json:"name"`
	Code               string          `json:"code"`
	FormStructure      []FormStructure `json:"form_structure"`
	ComponentStructure string          `json:"component_structure"`
	Version            int             `json:"version"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type FormInstanceBo struct {
	ID           int       `json:"id"`
	FormDefineID int          `json:"form_define_id"`
	FormData     []FormData   `json:"form_data"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
