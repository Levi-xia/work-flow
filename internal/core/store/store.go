package store

import (
	"encoding/json"
	"log"
	"workflow/internal/core/constants"
	"workflow/internal/core/model"
	"workflow/internal/serctx"

	_ "github.com/go-sql-driver/mysql"
)

// 包含初始化、取用逻辑
type store struct {
	processDefineStore   ProcessDefineStore   // 定义存储
	processInstanceStore ProcessInstanceStore // 实例存储
	processTaskStore     ProcessTaskStore     // 任务存储
}

// 全局store实例
var globalStore *store

func init() {
	log.Println("init store")
	globalStore = &store{
		processDefineStore:   &MySQLProcessDefineStore{},
		processInstanceStore: &MySQLProcessInstanceStore{},
		processTaskStore:     &MySQLProcessTaskStore{},
	}
}

// 获取定义存储
func GetProcessDefineStore() ProcessDefineStore {
	return globalStore.processDefineStore
}

// 获取实例存储
func GetProcessInstanceStore() ProcessInstanceStore {
	return globalStore.processInstanceStore
}

// 获取任务存储
func GetProcessTaskStore() ProcessTaskStore {
	return globalStore.processTaskStore
}

type MySQLProcessDefineStore struct{}

func (s *MySQLProcessDefineStore) GetProcessDefine(id int) (*model.ProcessDefineModel, error) {
	define := &model.ProcessDefineModel{}
	err := serctx.SerCtx.Db.Get(define, "SELECT * FROM process_define WHERE id = ?", id)
	return define, err
}
func (s *MySQLProcessDefineStore) CreateProcessDefine(meta *model.ProcessDefineModel) (int, error) {
	result, err := serctx.SerCtx.Db.NamedExec(`
		INSERT INTO process_define (code, user_id, name, version, content)
		VALUES (:code, :name, :version, :content)`, meta)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (s *MySQLProcessDefineStore) GetProcessDefinesByCode(code string) ([]*model.ProcessDefineModel, error) {
	var defines []*model.ProcessDefineModel
	err := serctx.SerCtx.Db.Select(&defines, "SELECT * FROM process_define WHERE code = ?", code)
	return defines, err
}

type MySQLProcessInstanceStore struct{}

func (s *MySQLProcessInstanceStore) CreateProcessInstance(meta *model.ProcessInstanceModel) (int, error) {
	result, err := serctx.SerCtx.Db.Exec(`
		INSERT INTO process_instance (process_define_id, user_id, status, variables)
		VALUES (?, ?, ?)`, meta.ProcessDefineID, meta.UserID, meta.Status, meta.Variables)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}
func (s *MySQLProcessInstanceStore) FinishProcessInstance(id int) error {
	_, err := serctx.SerCtx.Db.Exec("UPDATE process_instance SET status = ? WHERE id = ?", constants.PROCESSINSTANCESTATUSFINISH, id)
	return err
}
func (s *MySQLProcessInstanceStore) CancelProcessInstance(id int) error {
	_, err := serctx.SerCtx.Db.Exec("UPDATE process_instance SET status = ? WHERE id = ?", constants.PROCESSINSTANCESTATUSCANCEL, id)
	return err
}
func (s *MySQLProcessInstanceStore) GetProcessInstance(id int) (*model.ProcessInstanceModel, error) {
	instance := &model.ProcessInstanceModel{}
	if err := serctx.SerCtx.Db.Get(instance, "SELECT * FROM process_instance WHERE id = ?", id); err != nil {
		return nil, err
	}
	return instance, nil
}
func (s *MySQLProcessInstanceStore) UpdateProcessInstanceVariables(id int, variables map[string]interface{}) error {
	variablesBytes, err := json.Marshal(variables)
	if err != nil {
		return err
	}
	_, err = serctx.SerCtx.Db.Exec("UPDATE process_instance SET variables = ? WHERE id = ?", string(variablesBytes), id)
	return err
}

type MySQLProcessTaskStore struct{}

func (s *MySQLProcessTaskStore) CreateProcessTask(meta *model.ProcessTaskModel) (int, error) {
	result, err := serctx.SerCtx.Db.Exec(`
		INSERT INTO process_task (process_instance_id, form_instance_id, code, name, status, variables)
		VALUES (?, ?, ?, ?, ?, ?)`, meta.ProcessInstanceID, meta.FormInstanceID, meta.Code, meta.Name, meta.Status, meta.Variables)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}
func (s *MySQLProcessTaskStore) FinishProcessTask(id int, userId int, variables map[string]interface{}) error {
	variablesBytes, err := json.Marshal(variables)
	if err != nil {
		return err
	}
	_, err = serctx.SerCtx.Db.Exec("UPDATE process_task SET status = ?, user_id = ? variables = ? WHERE id = ?", constants.PROCESSTASKSTATUSFINISH, userId, string(variablesBytes), id)
	return err
}
func (s *MySQLProcessTaskStore) GetProcessTask(id int) (*model.ProcessTaskModel, error) {
	task := &model.ProcessTaskModel{}
	if err := serctx.SerCtx.Db.Get(task, "SELECT * FROM process_task WHERE id = ?", id); err != nil {
		return nil, err
	}
	return task, nil
}
func (s *MySQLProcessTaskStore) GetRunningTasks(instanceID int) ([]*model.ProcessTaskModel, error) {
	var tasks []*model.ProcessTaskModel
	err := serctx.SerCtx.Db.Select(&tasks, "SELECT * FROM process_task WHERE process_instance_id = ? AND status = ?", instanceID, constants.PROCESSTASKSTATUSDOING)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
