package store

import (
	"workflow/internal/form/model"
	"workflow/internal/serctx"
)

type store struct {
	formDefineStore   FormDefineStore   // 定义存储
	formInstanceStore FormInstanceStore // 实例存储
}

// 全局store实例
var GlobalStore *store

func init() {
	GlobalStore = &store{
		formDefineStore:   &MySQLFormDefineStore{},
		formInstanceStore: &MySQLFormInstanceStore{},
	}
}

// 获取定义存储
func GetFormDefineStore() FormDefineStore {
	return GlobalStore.formDefineStore
}

// 获取实例存储
func GetFormInstanceStore() FormInstanceStore {
	return GlobalStore.formInstanceStore
}

type MySQLFormDefineStore struct{}

func (s *MySQLFormDefineStore) CreateFormDefine(meta *model.FormDefineModel) (int, error) {
	result, err := serctx.SerCtx.Db.NamedExec(`
		INSERT INTO form_define (code, name, version, form_structure, component_structure)
		VALUES (:code, :name, :version, :form_structure, :component_structure)`, meta)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (s *MySQLFormDefineStore) GetFormDefine(id int) (*model.FormDefineModel, error) {
	define := &model.FormDefineModel{}
	err := serctx.SerCtx.Db.Get(define, "SELECT * FROM form_define WHERE id = ?", id)
	return define, err
}

func (s *MySQLFormDefineStore) GetFormDefinesByCode(code string) ([]*model.FormDefineModel, error) {
	var defines []*model.FormDefineModel
	err := serctx.SerCtx.Db.Select(&defines, "SELECT * FROM form_define WHERE code = ?", code)
	return defines, err
}

type MySQLFormInstanceStore struct{}

func (s *MySQLFormInstanceStore) CreateFormInstance(formInstance *model.FormInstanceModel) (int, error) {
	result, err := serctx.SerCtx.Db.NamedExec(`
		INSERT INTO form_instance (form_define_id, form_data)
		VALUES (:form_define_id, :form_data)`, formInstance)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (s *MySQLFormInstanceStore) GetFormInstance(id int) (*model.FormInstanceModel, error) {
	formInstance := &model.FormInstanceModel{}
	err := serctx.SerCtx.Db.Get(formInstance, "SELECT * FROM form_instance WHERE id = ?", id)
	return formInstance, err
}

func (s *MySQLFormInstanceStore) UpdateFormInstanceFormData(formInstanceID int, formData string) error {
	_, err := serctx.SerCtx.Db.Exec("UPDATE form_instance SET form_data = ? WHERE id = ?", formData, formInstanceID)
	return err
}
