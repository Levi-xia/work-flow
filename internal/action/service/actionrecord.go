package service

import (
	"encoding/json"
	"workflow/internal/action/bo"
	"workflow/internal/action/model"
	"workflow/internal/action/store"
)

type ActionRecord struct {
	Define *ActionDefine
	Meta   *bo.ActionRecordBo
	store  store.ActionRecordStore
}

func NewActionRecord(record *bo.ActionRecordBo) (*ActionRecord, error) {
	var (
		err               error
		inputBytes        []byte
		outputBytes       []byte
		actionRecord      *ActionRecord
		actionRecordModel *model.ActionRecordModel
	)
	if inputBytes, err = json.Marshal(record.Input); err != nil {
		return nil, err
	}
	if outputBytes, err = json.Marshal(record.Output); err != nil {
		return nil, err
	}
	actionRecordModel = &model.ActionRecordModel{
		ActionDefineID: record.ActionDefineID,
		ProcessTaskID:  record.ProcessTaskID,
		Input:          string(inputBytes),
		Output:         string(outputBytes),
	}
	actionRecord = &ActionRecord{
		Meta:  record,
		store: store.GetActionRecordStore(),
	}
	if actionRecord.Meta.ID, err = store.GetActionRecordStore().CreateActionRecord(actionRecordModel); err != nil {
		return nil, err
	}
	return actionRecord, nil
}

func UpdateActionRecordResult(id int, input map[string]interface{}, output map[string]interface{}) error {
	return store.GetActionRecordStore().UpdateActionRecordResult(id, input, output)
}
