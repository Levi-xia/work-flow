package service

import (
	"encoding/json"
	"sort"
	"workflow/internal/action/bo"
	"workflow/internal/action/model"
	"workflow/internal/action/store"
)

type ActionDefine struct {
	Meta  *bo.ActionDefineBo
	store store.ActionDefineStore
}

// 创建ActionDefine
func NewActionDefine(name, code, protocol string, content *bo.ActionContent,
	inputStructs []*bo.ParamsStruct, outputChecks []*bo.OutputCheck) (*ActionDefine, error) {

	var (
		err               error
		defineModel       *model.ActionDefineModel
		lastVersion       *model.ActionDefineModel
		defineMeta        *bo.ActionDefineBo
		contentBytes      []byte
		inputStructsBytes []byte
		outputChecksBytes []byte
	)
	if contentBytes, err = json.Marshal(content); err != nil {
		return nil, err
	}
	if inputStructsBytes, err = json.Marshal(inputStructs); err != nil {
		return nil, err
	}
	if outputChecksBytes, err = json.Marshal(outputChecks); err != nil {
		return nil, err
	}
	defineModel = &model.ActionDefineModel{
		Name:         name,
		Code:         code,
		Protocol:     protocol,
		Content:      string(contentBytes),
		Version:      1,
		InputStructs: string(inputStructsBytes),
		OutputChecks: string(outputChecksBytes),
	}
	if lastVersion, err = getLatestVersion(code); err != nil {
		return nil, err
	}
	if lastVersion != nil {
		defineModel.Version = lastVersion.Version + 1
	}
	// 写入数据库
	if defineModel.ID, err = store.GetActionDefineStore().CreateActionDefine(defineModel); err != nil {
		return nil, err
	}
	if defineMeta, err = defineModel.ToBo(); err != nil {
		return nil, err
	}
	return &ActionDefine{
		Meta:  defineMeta,
		store: store.GetActionDefineStore(),
	}, nil
}

// 获取ActionDefine
func GetActionDefine(id int) (*bo.ActionDefineBo, error) {
	defineModel, err := store.GetActionDefineStore().GetActionDefine(id)
	if err != nil {
		return nil, err
	}
	return defineModel.ToBo()
}

// 获取最新版本
func getLatestVersion(code string) (*model.ActionDefineModel, error) {
	defineModels, err := store.GetActionDefineStore().GetActionDefinesByCode(code)
	if err != nil {
		return nil, err
	}
	if len(defineModels) == 0 {
		return nil, nil
	}
	// 根据version排序
	sort.Slice(defineModels, func(i, j int) bool {
		return defineModels[i].Version > defineModels[j].Version
	})
	return defineModels[0], nil
}
