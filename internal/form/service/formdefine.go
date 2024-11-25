package service

import (
	"sort"
	"workflow/internal/form/bo"
	"workflow/internal/form/model"
	"workflow/internal/form/store"
	"workflow/internal/utils"
)

type FormDefine struct {
	Meta  *model.FormDefineModel
	store store.FormDefineStore
}

func NewFormDefine(name, code, formStructure, componentStructure string) (*FormDefine, error) {
	latestVersion, err := getLatestVersion(code)
	if err != nil {
		return nil, err
	}
	define := &FormDefine{
		Meta: &model.FormDefineModel{
			Name:               name,
			Code:               code,
			Version:            1,
			FormStructure:      formStructure,
			ComponentStructure: componentStructure,
		},
		store: store.GetFormDefineStore(),
	}
	if latestVersion != nil {
		define.Meta.Version = latestVersion.Version + 1
	}
	// 压缩JSON
	if define.Meta.FormStructure, err = utils.CompactJSON(formStructure); err != nil {
		return nil, err
	}
	if define.Meta.ComponentStructure, err = utils.CompactJSON(componentStructure); err != nil {
		return nil, err
	}
	// 写入数据库
	if define.Meta.ID, err = store.GetFormDefineStore().CreateFormDefine(define.Meta); err != nil {
		return nil, err
	}
	return define, nil
}

func GetFormDefine(id int) (*bo.FormDefineBo, error) {
	formDefineModel, err := store.GetFormDefineStore().GetFormDefine(id)
	if err != nil {
		return nil, err
	}
	return formDefineModel.ToBo()
}

func getLatestVersion(code string) (*model.FormDefineModel, error) {
	formDefineModels, err := store.GetFormDefineStore().GetFormDefinesByCode(code)
	if err != nil {
		return nil, err
	}
	if len(formDefineModels) == 0 {
		return nil, nil
	}
	sort.Slice(formDefineModels, func(i, j int) bool {
		return formDefineModels[i].Version > formDefineModels[j].Version
	})
	return formDefineModels[0], nil
}
