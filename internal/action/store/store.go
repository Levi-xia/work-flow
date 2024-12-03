package store

import (
	"log"
	"workflow/internal/action/model"
	"workflow/internal/serctx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// 包含初始化、取用逻辑
type store struct {
	actionDefineStore ActionDefineStore // 定义存储
	actionRecordStore ActionRecordStore // 记录存储
}

// 全局store实例
var globalStore *store

func init() {
	log.Println("init action store")
	globalStore = &store{
		actionDefineStore: &MySQLActionDefineStore{},
		actionRecordStore: &MySQLActionRecordStore{},
	}
}

// 获取定义存储
func GetActionDefineStore() ActionDefineStore {
	return globalStore.actionDefineStore
}

// 获取实例存储
func GetActionRecordStore() ActionRecordStore {
	return globalStore.actionRecordStore
}

type MySQLActionDefineStore struct{}

func (s *MySQLActionDefineStore) CreateActionDefine(defineModel *model.ActionDefineModel) (int, error) {
	result, err := serctx.SerCtx.Db.NamedExec(`
		INSERT INTO action_define (code, name, user_id, version, protocol, content, input_structs, output_checks)
		VALUES (:code, :name, :user_id, :version, :protocol, :content, :input_structs, :output_checks)`, defineModel)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (s *MySQLActionDefineStore) GetActionDefine(id int) (*model.ActionDefineModel, error) {
	var defineModel model.ActionDefineModel
	err := serctx.SerCtx.Db.Get(&defineModel, "SELECT * FROM action_define WHERE id = ?", id)
	return &defineModel, err
}

func (s *MySQLActionDefineStore) GetActionDefines(ids []int) ([]*model.ActionDefineModel, error) {
	var defineModels []*model.ActionDefineModel
	query, args, err := sqlx.In("SELECT * FROM action_define WHERE id in (?)", ids)
	if err != nil {
		return nil, err
	}
	err = serctx.SerCtx.Db.Select(&defineModels, query, args...)
	return defineModels, err
}

func (s *MySQLActionDefineStore) GetActionDefinesByCode(code string) ([]*model.ActionDefineModel, error) {
	var defineModels []*model.ActionDefineModel
	err := serctx.SerCtx.Db.Select(&defineModels, "SELECT * FROM action_define WHERE code = ?", code)
	return defineModels, err
}

type MySQLActionRecordStore struct{}

func (s *MySQLActionRecordStore) CreateActionRecord(record *model.ActionRecordModel) (int, error) {
	result, err := serctx.SerCtx.Db.NamedExec(`
		INSERT INTO action_record (action_define_id, process_task_id, input, output)
		VALUES (:action_define_id, :process_task_id, :input, :output)`, record)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (s *MySQLActionRecordStore) GetActionRecord(id int) (*model.ActionRecordModel, error) {
	var record model.ActionRecordModel
	err := serctx.SerCtx.Db.Get(&record, "SELECT * FROM action_record WHERE id = ?", id)
	return &record, err
}

func (s *MySQLActionRecordStore) GetActionRecordsByProcessTaskID(processTaskID int) ([]*model.ActionRecordModel, error) {
	var records []*model.ActionRecordModel
	err := serctx.SerCtx.Db.Select(&records, "SELECT * FROM action_record WHERE process_task_id = ?", processTaskID)
	return records, err
}

func (s *MySQLActionRecordStore) UpdateActionRecordResult(id int, input map[string]interface{}, output map[string]interface{}) error {
	_, err := serctx.SerCtx.Db.NamedExec(`
		UPDATE action_record SET input = :input, output = :output WHERE id = :id`, map[string]interface{}{
		"id":     id,
		"input":  input,
		"output": output,
	})
	return err
}
