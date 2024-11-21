package service

import (
	"sort"
	"workflow/internal/core/bo"
	"workflow/internal/core/model"
	"workflow/internal/core/parser"
	"workflow/internal/core/process"
	"workflow/internal/core/store"
	"workflow/internal/utils"
)

type ProcessDefine struct {
	Meta  *bo.ProcessDefineBo
	store store.ProcessDefineStore
}

// NewProcessDefine 用于创建一个新的流程定义
func NewProcessDefine(content string) (*ProcessDefine, error) {
	var (
		err         error
		process     *process.Process
		define      *ProcessDefine
		defineModel *model.ProcessDefineModel
		defineBo    *bo.ProcessDefineBo
	)
	// 解析
	if process, err = parser.Parser2Process(content); err != nil {
		return nil, err
	}
	// 获取最新版本
	processDefineModel, err := getLatestVersion(process.Code)
	if err != nil {
		return nil, err
	}
	// 压缩JSON
	compactContent, err := utils.CompactJSON(content)
	if err != nil {
		return nil, err
	}
	defineModel = &model.ProcessDefineModel{
		Code:    process.Code,
		Name:    process.Name,
		Content: compactContent,
		Version: 1,
	}
	if processDefineModel != nil {
		defineModel.Version = processDefineModel.Version + 1
	}
	// 转换为BO
	if defineBo, err = defineModel.ToBo(); err != nil {
		return nil, err
	}
	define = &ProcessDefine{
		Meta:  defineBo,
		store: store.GetProcessDefineStore(),
	}
	// 写入数据库
	if define.Meta.ID, err = define.store.CreateProcessDefine(defineModel); err != nil {
		return nil, err
	}
	return define, nil
}

func GetProcessDefine(processDefineId int) (*bo.ProcessDefineBo, error) {
	defineModel, err := store.GetProcessDefineStore().GetProcessDefine(processDefineId)
	if err != nil {
		return nil, err
	}
	return defineModel.ToBo()
}

// 获取最高版本的定义
func getLatestVersion(code string) (*model.ProcessDefineModel, error) {
	processDefineModels, err := store.GetProcessDefineStore().GetProcessDefinesByCode(code)
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
