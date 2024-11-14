package base

import (
	"time"
	"workflow/internal/core/constants"
)

type HttpAction struct {
	Url     string
	Method  string
	Timeout time.Duration
	Headers map[string]string
}

type Action struct {
	ActionType constants.ActionType
	HttpAction HttpAction
	Params     map[string]interface{}
}