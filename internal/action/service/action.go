package service

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	"workflow/internal/action/bo"
	"workflow/internal/action/constants"
	"workflow/internal/utils"
)

// 执行actions
func ExecuteActions(actionDefineIds []int, processTaskID int, params map[string]interface{}) error {
	if params == nil {
		params = make(map[string]interface{})
	}
	if len(actionDefineIds) == 0 {
		return nil
	}
	actionDefines, err := GetActionDefines(actionDefineIds)
	if err != nil {
		return err
	}
	// check入参
	if err := checkParams(actionDefines, params); err != nil {
		return err
	}
	// 执行action
	if err := execActions(actionDefines, processTaskID, params); err != nil {
		return err
	}
	return nil
}

// 检查入参
func checkParams(actionDefines []*bo.ActionDefineBo, params map[string]interface{}) error {
	wg := sync.WaitGroup{}
	errChan := make(chan error)

	for _, actionDefine := range actionDefines {
		wg.Add(1)
		go func(actionDefine *bo.ActionDefineBo) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic in action check params: %v", r)
				}
			}()
			for _, inputStruct := range actionDefine.InputStructs {
				if err := checkParam(&inputStruct, params); err != nil {
					errChan <- err
				}
			}
		}(actionDefine)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		return err
	}
	return nil
}

// 执行action
func execActions(actionDefines []*bo.ActionDefineBo, processTaskID int, params map[string]interface{}) error {
	// 并发执行action
	wg := sync.WaitGroup{}
	errChan := make(chan error)
	for _, actionDefine := range actionDefines {
		wg.Add(1)
		go func(actionDefine *bo.ActionDefineBo) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic in action exec: %v", r)
				}
			}()
			if _, err := execAction(actionDefine, processTaskID, params); err != nil {
				errChan <- err
			}
		}(actionDefine)
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		return err
	}
	return nil
}

// 执行单个action
func execAction(actionDefine *bo.ActionDefineBo, processTaskID int, params map[string]interface{}) (interface{}, error) {
	var (
		err          error
		result       []byte
		actionParams map[string]interface{} = make(map[string]interface{})
		resultMap    map[string]interface{} = make(map[string]interface{})
	)
	// 成功失败均要创建action记录
	defer func() {
		if err != nil {
			resultMap = map[string]interface{}{
				"error": err.Error(),
			}
		}
		record := &bo.ActionRecordBo{
			ActionDefineID: actionDefine.ID,
			ProcessTaskID:  processTaskID,
			Input:          actionParams,
			Output:         resultMap,
		}
		if _, err := NewActionRecord(record); err != nil {
			log.Println("create action record failed: ", err)
		}
	}()
	if actionDefine == nil {
		return nil, fmt.Errorf("action define is nil")
	}
	// 根据action定义的inputStruct从params中获取参数
	for _, inputStruct := range actionDefine.InputStructs {
		if value, exists := params[inputStruct.Key]; exists {
			actionParams[inputStruct.Key] = value
		}
	}
	// 执行action
	switch actionDefine.Protocol {
	case constants.ActionProtocolHttp:
		httpAction := actionDefine.Content.HttpAction
		if result, err = utils.HttpDo(httpAction.Url, string(httpAction.Method), actionParams,
			utils.WithHeaders(httpAction.Headers),
			utils.WithTimeout(time.Duration(httpAction.Timeout)*time.Millisecond)); err != nil {
			return result, err
		}
	default:
		return nil, fmt.Errorf("unsupported action protocol %s", actionDefine.Protocol)
	}
	if err := json.Unmarshal(result, &resultMap); err != nil {
		return result, err
	}
	// 对执行结果进行检查
	if err = checkOutput(actionDefine.OutputChecks, resultMap); err != nil {
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
