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
	}
	
	err := ExecuteActions([]int{100002,100003,100004,100005}, 99999, params)
	if err != nil {
		t.Error(err)
	}
	t.Log(params)
}

func TestCreateActionDefine(t *testing.T) {
	serctx.InitServerContext()
	define, err := NewActionDefine("测试", "test", "http", &bo.ActionContent{
		HttpAction: bo.HttpAction{
			Url:     "http://127.0.0.1:8080/workflow/action/v1/sendSms",
			Method:  constants.HttpMethodGet,
			Headers: make(map[string]string),
		},
	}, []*bo.ParamsStruct{
		{
			Type:     constants.ParamsStructTypeString,
			Key:      "key1",
			Required: true,
		},
		{
			Type:     constants.ParamsStructTypeInt,
			Key:      "key2",
			Required: true,
		},
		{
			Type:     constants.ParamsStructTypeString,
			Key:      "key3",
			Required: false,
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
