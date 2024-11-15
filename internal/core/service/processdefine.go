package service

import (
	"sort"
	"workflow/internal/core/model"
	"workflow/internal/core/parser"
	"workflow/internal/core/store"
	"workflow/internal/utils"
)

type ProcessDefine struct {
	Meta  *model.ProcessDefineModel
	Store store.ProcessDefineStore
}

// 创建流程定义
func NewProcessDefine(content string) (*ProcessDefine, error) {
	// 解析
	process, err := parser.Parser2Process(content)
	if err != nil {
		return nil, err
	}
	code := process.Code
	name := process.Name

	processDefineModel, err := getLatestVersion(code, store.GetProcessDefineStore())
	if err != nil {
		return nil, err
	}
	define := &ProcessDefine{
		Meta: &model.ProcessDefineModel{
			Code:    code,
			Name:    name,
			Content: content,
			Version: 1,
		},
		Store: store.GetProcessDefineStore(),
	}
	if processDefineModel != nil {
		define.Meta.Version = processDefineModel.Version + 1
	}
	// 压缩JSON
	compactContent, err := utils.CompactJSON(content)
	if err != nil {
		return nil, err
	}
	define.Meta.Content = compactContent
	// 写入数据库
	if define.Meta.ID, err = define.Store.CreateProcessDefine(define.Meta); err != nil {
		return nil, err
	}
	return define, nil
}

// 获取最高版本的定义
func getLatestVersion(code string, store store.ProcessDefineStore) (*model.ProcessDefineModel, error) {
	processDefineModels, err := store.GetProcessDefinesByCode(code)
	if err != nil {
		return nil, err
	}
	if len(processDefineModels) == 0 {
		return nil, nil
	}
	// 根据version排序
	sort.Slice(processDefineModels, func(i, j int) bool {
		return processDefineModels[i].Version > processDefineModels[j].Version
	})
	return processDefineModels[0], nil
}
