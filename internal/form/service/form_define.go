package service

import (
	"sort"
	"workflow/internal/form/model"
	"workflow/internal/form/store"
	"workflow/internal/utils"
)

type FormDefine struct {
	Meta  *model.FormDefineModel
	Store store.FormDefineStore
}

func NewFormDefine(code, formStructure, componentStructure string, store store.FormDefineStore) (*FormDefine, error) {
	latestVersion, err := getLatestVersion(code, store)
	if err != nil {
		return nil, err
	}
	define := &FormDefine{
		Meta: &model.FormDefineModel{
			Code:               code,
			Version:            1,
			FormStructure:      formStructure,
			ComponentStructure: componentStructure,
		},
		Store: store,
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
	if define.Meta.ID, err = define.Store.CreateFormDefine(define.Meta); err != nil {
		return nil, err
	}
	return define, nil
}

func getLatestVersion(code string, store store.FormDefineStore) (*model.FormDefineModel, error) {
	formDefineModels, err := store.GetFormDefinesByCode(code)
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
