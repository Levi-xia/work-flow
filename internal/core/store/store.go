package store

import (
	"encoding/json"
	"log"
	"workflow/internal/core/constants"
	"workflow/internal/core/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// 包含初始化、取用逻辑
type store struct {
	processDefineStore   ProcessDefineStore   // 定义存储
	processInstanceStore ProcessInstanceStore // 实例存储
	processTaskStore     ProcessTaskStore     // 任务存储
}

// 全局store实例
var GlobalStore *store

var db *sqlx.DB

func init() {
	log.Println("init store")
	GlobalStore = &store{
		processDefineStore:   &MySQLProcessDefineStore{},
		processInstanceStore: &MySQLProcessInstanceStore{},
		processTaskStore:     &MySQLProcessTaskStore{},
	}
	// 初始化数据库连接
	dsn := "root:Mysqlxwy9264@tcp(127.0.0.1:3306)/workflow?charset=utf8mb4&parseTime=True"
	db = sqlx.MustConnect("mysql", dsn)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
}

// 获取定义存储
func GetProcessDefineStore() ProcessDefineStore {
	return GlobalStore.processDefineStore
}

// 获取实例存储
func GetProcessInstanceStore() ProcessInstanceStore {
	return GlobalStore.processInstanceStore
}

// 获取任务存储
func GetProcessTaskStore() ProcessTaskStore {
	return GlobalStore.processTaskStore
}

type MySQLProcessDefineStore struct{}

func (s *MySQLProcessDefineStore) GetProcessDefine(id int) (*model.ProcessDefineModel, error) {
	define := &model.ProcessDefineModel{}
	err := db.Get(define, "SELECT * FROM process_define WHERE id = ?", id)
	return define, err
}
func (s *MySQLProcessDefineStore) CreateProcessDefine(meta *model.ProcessDefineModel) (int, error) {
	result, err := db.NamedExec(`
		INSERT INTO process_define (code, name, version, content)
		VALUES (:code, :name, :version, :content)`, meta)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (s *MySQLProcessDefineStore) GetProcessDefinesByCode(code string) ([]*model.ProcessDefineModel, error) {
	var defines []*model.ProcessDefineModel
	err := db.Select(&defines, "SELECT * FROM process_define WHERE code = ?", code)
	return defines, err
}

type MySQLProcessInstanceStore struct{}

func (s *MySQLProcessInstanceStore) CreateProcessInstance(meta *model.ProcessInstanceModel) (int, error) {
	variables, err := json.Marshal(meta.Variables)
	if err != nil {
		return 0, err
	}
	result, err := db.Exec(`
		INSERT INTO process_instance (process_define_id, status, variables)
		VALUES (?, ?, ?)`, meta.ProcessDefineID, meta.Status, string(variables))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}
func (s *MySQLProcessInstanceStore) FinishProcessInstance(id int) error {
	_, err := db.Exec("UPDATE process_instance SET status = ? WHERE id = ?", constants.PROCESSINSTANCESTATUSFINISH, id)
	return err
}
func (s *MySQLProcessInstanceStore) CancelProcessInstance(id int) error {
	_, err := db.Exec("UPDATE process_instance SET status = ? WHERE id = ?", constants.PROCESSINSTANCESTATUSCANCEL, id)
	return err
}
func (s *MySQLProcessInstanceStore) GetProcessInstance(id int) (*model.ProcessInstanceModel, error) {
	// 创建一个临时结构体来接收原始数据
	type tempInstance struct {
		*model.ProcessInstanceModel
		Variables string `db:"variables"`
	}
	tmp := &tempInstance{ProcessInstanceModel: &model.ProcessInstanceModel{}}
	err := db.Get(tmp, "SELECT * FROM process_instance WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	// 将 JSON 字符串解析为 map
	if tmp.Variables != "" {
		if err := json.Unmarshal([]byte(tmp.Variables), &tmp.ProcessInstanceModel.Variables); err != nil {
			return nil, err
		}
	}
	return tmp.ProcessInstanceModel, nil
}
func (s *MySQLProcessInstanceStore) AddVariables(id int, k string, v interface{}) error {
	_, err := db.Exec(`
		UPDATE process_instance 
		SET variables = JSON_SET(COALESCE(variables, '{}'), ?, CAST(? AS JSON))
		WHERE id = ?`, "$."+k, v, id)
	return err
}
func (s *MySQLProcessInstanceStore) DeleteVariables(id int, variables ...string) error {
	query := `UPDATE process_instance SET variables = JSON_REMOVE(variables`
	for _, k := range variables {
		query += `, '$."` + k + `"'`
	}
	query += `), updated_at = NOW() WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

type MySQLProcessTaskStore struct{}

func (s *MySQLProcessTaskStore) CreateProcessTask(meta *model.ProcessTaskModel) (int, error) {
	variables, err := json.Marshal(meta.Variables)
	if err != nil {
		return 0, err
	}
	result, err := db.Exec(`
		INSERT INTO process_task (process_instance_id, code, name, status, variables)
		VALUES (?, ?, ?, ?, ?)`, meta.ProcessInstanceID, meta.Code, meta.Name, meta.Status, string(variables))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}
func (s *MySQLProcessTaskStore) FinishProcessTask(id int) error {
	_, err := db.Exec("UPDATE process_task SET status = ? WHERE id = ?", constants.PROCESSTASKSTATUSFINISH, id)
	return err
}
func (s *MySQLProcessTaskStore) GetProcessTask(id int) (*model.ProcessTaskModel, error) {
	// 创建一个临时结构体来接收原始数据
	type tempTask struct {
		*model.ProcessTaskModel
		Variables string `db:"variables"`
	}
	tmp := &tempTask{ProcessTaskModel: &model.ProcessTaskModel{}}
	if err := db.Get(tmp, "SELECT * FROM process_task WHERE id = ?", id); err != nil {
		return nil, err
	}
	// 将 JSON 字符串解析为 map
	if tmp.Variables != "" {
		if err := json.Unmarshal([]byte(tmp.Variables), &tmp.ProcessTaskModel.Variables); err != nil {
			return nil, err
		}
	}
	return tmp.ProcessTaskModel, nil
}
func (s *MySQLProcessTaskStore) GetRunningTasks(instanceID int) ([]*model.ProcessTaskModel, error) {
	// 创建一个临时结构体来接收原始数据
	type tempTask struct {
		*model.ProcessTaskModel
		Variables string `db:"variables"`
	}
	var tasks []*tempTask
	err := db.Select(&tasks, "SELECT * FROM process_task WHERE process_instance_id = ? AND status = ?", instanceID, constants.PROCESSTASKSTATUSDOING)
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		if task.Variables != "" {
			if err := json.Unmarshal([]byte(task.Variables), &task.ProcessTaskModel.Variables); err != nil {
				return nil, err
			}
		}
	}
	var tasks2 []*model.ProcessTaskModel
	for _, task := range tasks {
		tasks2 = append(tasks2, task.ProcessTaskModel)
	}
	return tasks2, nil
}
