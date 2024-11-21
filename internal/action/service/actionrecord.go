package service

import (
	"workflow/internal/action/bo"
	"workflow/internal/action/store"
)

type ActionRecord struct {
	Define *ActionDefine
	Meta   *bo.ActionRecordBo
	store  store.ActionRecordStore
}


func UpdateActionRecordResult(id int, input map[string]interface{}, output map[string]interface{}) error {
	return store.GetActionRecordStore().UpdateActionRecordResult(id, input, output)
}
