package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"workflow/internal/action/bo"
	"workflow/internal/action/constants"
	"workflow/internal/utils"
)

// 执行actions
func ExecuteAction(actions []*ActionRecord, params map[string]interface{}) error {
	if params == nil {
		params = make(map[string]interface{})
	}
	if len(actions) == 0 {
		return nil
	}
	// check入参
	if err := checkParams(actions, params); err != nil {
		return err
	}
	// 执行action
	if err := execActions(actions, params); err != nil {
		return err
	}
	return nil
}

// 检查入参
func checkParams(actions []*ActionRecord, params map[string]interface{}) error {
	wg := sync.WaitGroup{}
	errChan := make(chan error)

	for _, action := range actions {
		wg.Add(1)
		go func(action *ActionRecord) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic in action check params: %v", r)
				}
			}()
			for _, inputStruct := range action.Define.Meta.InputStructs {
				if err := checkParam(&inputStruct, params); err != nil {
					errChan <- err
					break
				}
			}
		}(action)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		return err
	}
	return nil
}

// 执行action
func execActions(actions []*ActionRecord, params map[string]interface{}) error {
	// 并发执行action
	wg := sync.WaitGroup{}
	errChan := make(chan error)
	for _, action := range actions {
		wg.Add(1)
		go func(action *ActionRecord) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic in action exec: %v", r)
				}
			}()
			if _, err := execAction(action, params); err != nil {
				errChan <- err
			}
		}(action)
	}
	return nil
}

// 执行单个action
func execAction(action *ActionRecord, params map[string]interface{}) (interface{}, error) {
	if action.Define == nil {
		return nil, fmt.Errorf("action define is nil")
	}
	var (
		err    error
		result []byte
	)
	// 根据action定义的inputStruct从params中获取参数
	actionParams := make(map[string]interface{})
	for _, inputStruct := range action.Define.Meta.InputStructs {
		if value, exists := params[inputStruct.Key]; exists {
			actionParams[inputStruct.Key] = value
		}
	}
	// 执行action
	switch action.Define.Meta.Protocol {
	case constants.ActionProtocolHttp:
		httpAction := action.Define.Meta.Content.HttpAction
		if result, err = utils.HttpDo(httpAction.Url, string(httpAction.Method), actionParams,
			utils.WithHeaders(httpAction.Headers),
			utils.WithTimeout(time.Duration(httpAction.Timeout)*time.Millisecond)); err != nil {
			return result, err
		}
	default:
		return nil, fmt.Errorf("unsupported action protocol %s", action.Define.Meta.Protocol)
	}
	resultMap := make(map[string]interface{})
	if err := json.Unmarshal(result, &resultMap); err != nil {
		return result, err
	}
	// 记录执行结果
	if err = UpdateActionRecordResult(action.Meta.ID, actionParams, resultMap); err != nil {
		return result, err
	}
	// 对执行结果进行检查
	if err = checkOutput(action.Define.Meta.OutputChecks, resultMap); err != nil {
		return result, err
	}
	return nil, nil
}

// 检查单个参数
func checkParam(inputStruct *bo.ParamsStruct, params map[string]interface{}) error {
	paramValue, exists := params[inputStruct.Key]
	// 检查必填参数是否存在
	if inputStruct.Required && !exists {
		return fmt.Errorf("required parameter %s is missing", inputStruct.Key)
	}
	// 如果参数存在,检查类型
	if exists {
		switch inputStruct.Type {
		case constants.ParamsStructTypeString:
			if _, ok := paramValue.(string); !ok {
				return fmt.Errorf("parameter %s should be string type", inputStruct.Key)
			}
		case constants.ParamsStructTypeInt:
			switch paramValue.(type) {
			case int, int32, int64:
			default:
				return fmt.Errorf("parameter %s should be int type", inputStruct.Key)
			}
		case constants.ParamsStructTypeFloat:
			switch paramValue.(type) {
			case float32, float64:
			default:
				return fmt.Errorf("parameter %s should be float type", inputStruct.Key)
			}
		case constants.ParamsStructTypeBool:
			if _, ok := paramValue.(bool); !ok {
				return fmt.Errorf("parameter %s should be boolean type", inputStruct.Key)
			}
		case constants.ParamsStructTypeArray:
			if _, ok := paramValue.([]interface{}); !ok {
				return fmt.Errorf("parameter %s should be array type", inputStruct.Key)
			}
		case constants.ParamsStructTypeObject:
			if _, ok := paramValue.(map[string]interface{}); !ok {
				return fmt.Errorf("parameter %s should be object type", inputStruct.Key)
			}
		default:
			return fmt.Errorf("unsupported parameter type %s", inputStruct.Type)
		}
	}
	return nil
}

// 检查输出
func checkOutput(outputChecks []bo.OutputCheck, resultMap map[string]interface{}) error {
	for _, outputCheck := range outputChecks {
		if value, exists := resultMap[outputCheck.Key]; !exists || value != outputCheck.Value {
			return fmt.Errorf("output check failed for key %s", outputCheck.Key)
		}
	}
	return nil
}
