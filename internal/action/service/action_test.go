package service

import (
	"testing"
	"workflow/internal/action/bo"
	"workflow/internal/action/constants"
	"workflow/internal/serctx"
)

func TestExecuteActions(t *testing.T) {
	serctx.InitServerContext()
	params := map[string]interface{}{
		"key1": "123",
		"key2": 123,
		"key3": "456",
		"key4": "789",
	}
	err := ExecuteActions([]int{100002}, 99999, params)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateActionDefine(t *testing.T) {
	serctx.InitServerContext()
	define, err := NewActionDefine(1, "通用请假action", "leave", "http", &bo.ActionContent{
		HttpAction: bo.HttpAction{
			Url:     "http://127.0.0.1:8080/workflow/action/v1/sendSms",
			Method:  constants.HttpMethodGet,
			Headers: make(map[string]string),
		},
	}, []*bo.ParamsStruct{
		{
			Type:     constants.ParamsStructTypeNumber,
			Key:      "days",
			Required: true,
		},
		{
			Type:     constants.ParamsStructTypeString,
			Key:      "reason",
			Required: true,
		},
	}, []*bo.OutputCheck{
		{
			Key:   "assert",
			Value: true,
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(define)
}
