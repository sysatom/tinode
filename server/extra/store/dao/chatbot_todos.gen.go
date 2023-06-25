// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/tinode/chat/server/extra/store/model"
)

func newTodo(db *gorm.DB, opts ...gen.DOOption) todo {
	_todo := todo{}

	_todo.todoDo.UseDB(db, opts...)
	_todo.todoDo.UseModel(&model.Todo{})

	tableName := _todo.todoDo.TableName()
	_todo.ALL = field.NewAsterisk(tableName)
	_todo.ID = field.NewInt32(tableName, "id")
	_todo.UID = field.NewString(tableName, "uid")
	_todo.Topic = field.NewString(tableName, "topic")
	_todo.Sequence = field.NewInt64(tableName, "sequence")
	_todo.Content = field.NewString(tableName, "content")
	_todo.Category = field.NewString(tableName, "category")
	_todo.Remark = field.NewString(tableName, "remark")
	_todo.Priority = field.NewInt64(tableName, "priority")
	_todo.IsRemindAtTime = field.NewInt32(tableName, "is_remind_at_time")
	_todo.RemindAt = field.NewInt64(tableName, "remind_at")
	_todo.RepeatMethod = field.NewString(tableName, "repeat_method")
	_todo.RepeatRule = field.NewString(tableName, "repeat_rule")
	_todo.RepeatEndAt = field.NewInt64(tableName, "repeat_end_at")
	_todo.Complete = field.NewInt32(tableName, "complete")
	_todo.CreatedAt = field.NewTime(tableName, "created_at")
	_todo.UpdatedAt = field.NewTime(tableName, "updated_at")

	_todo.fillFieldMap()

	return _todo
}

type todo struct {
	todoDo

	ALL            field.Asterisk
	ID             field.Int32
	UID            field.String
	Topic          field.String
	Sequence       field.Int64
	Content        field.String
	Category       field.String
	Remark         field.String
	Priority       field.Int64
	IsRemindAtTime field.Int32
	RemindAt       field.Int64
	RepeatMethod   field.String
	RepeatRule     field.String
	RepeatEndAt    field.Int64
	Complete       field.Int32
	CreatedAt      field.Time
	UpdatedAt      field.Time

	fieldMap map[string]field.Expr
}

func (t todo) Table(newTableName string) *todo {
	t.todoDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t todo) As(alias string) *todo {
	t.todoDo.DO = *(t.todoDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *todo) updateTableName(table string) *todo {
	t.ALL = field.NewAsterisk(table)
	t.ID = field.NewInt32(table, "id")
	t.UID = field.NewString(table, "uid")
	t.Topic = field.NewString(table, "topic")
	t.Sequence = field.NewInt64(table, "sequence")
	t.Content = field.NewString(table, "content")
	t.Category = field.NewString(table, "category")
	t.Remark = field.NewString(table, "remark")
	t.Priority = field.NewInt64(table, "priority")
	t.IsRemindAtTime = field.NewInt32(table, "is_remind_at_time")
	t.RemindAt = field.NewInt64(table, "remind_at")
	t.RepeatMethod = field.NewString(table, "repeat_method")
	t.RepeatRule = field.NewString(table, "repeat_rule")
	t.RepeatEndAt = field.NewInt64(table, "repeat_end_at")
	t.Complete = field.NewInt32(table, "complete")
	t.CreatedAt = field.NewTime(table, "created_at")
	t.UpdatedAt = field.NewTime(table, "updated_at")

	t.fillFieldMap()

	return t
}

func (t *todo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *todo) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 16)
	t.fieldMap["id"] = t.ID
	t.fieldMap["uid"] = t.UID
	t.fieldMap["topic"] = t.Topic
	t.fieldMap["sequence"] = t.Sequence
	t.fieldMap["content"] = t.Content
	t.fieldMap["category"] = t.Category
	t.fieldMap["remark"] = t.Remark
	t.fieldMap["priority"] = t.Priority
	t.fieldMap["is_remind_at_time"] = t.IsRemindAtTime
	t.fieldMap["remind_at"] = t.RemindAt
	t.fieldMap["repeat_method"] = t.RepeatMethod
	t.fieldMap["repeat_rule"] = t.RepeatRule
	t.fieldMap["repeat_end_at"] = t.RepeatEndAt
	t.fieldMap["complete"] = t.Complete
	t.fieldMap["created_at"] = t.CreatedAt
	t.fieldMap["updated_at"] = t.UpdatedAt
}

func (t todo) clone(db *gorm.DB) todo {
	t.todoDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t todo) replaceDB(db *gorm.DB) todo {
	t.todoDo.ReplaceDB(db)
	return t
}

type todoDo struct{ gen.DO }

func (t todoDo) Debug() *todoDo {
	return t.withDO(t.DO.Debug())
}

func (t todoDo) WithContext(ctx context.Context) *todoDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t todoDo) ReadDB() *todoDo {
	return t.Clauses(dbresolver.Read)
}

func (t todoDo) WriteDB() *todoDo {
	return t.Clauses(dbresolver.Write)
}

func (t todoDo) Session(config *gorm.Session) *todoDo {
	return t.withDO(t.DO.Session(config))
}

func (t todoDo) Clauses(conds ...clause.Expression) *todoDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t todoDo) Returning(value interface{}, columns ...string) *todoDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t todoDo) Not(conds ...gen.Condition) *todoDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t todoDo) Or(conds ...gen.Condition) *todoDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t todoDo) Select(conds ...field.Expr) *todoDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t todoDo) Where(conds ...gen.Condition) *todoDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t todoDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *todoDo {
	return t.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (t todoDo) Order(conds ...field.Expr) *todoDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t todoDo) Distinct(cols ...field.Expr) *todoDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t todoDo) Omit(cols ...field.Expr) *todoDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t todoDo) Join(table schema.Tabler, on ...field.Expr) *todoDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t todoDo) LeftJoin(table schema.Tabler, on ...field.Expr) *todoDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t todoDo) RightJoin(table schema.Tabler, on ...field.Expr) *todoDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t todoDo) Group(cols ...field.Expr) *todoDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t todoDo) Having(conds ...gen.Condition) *todoDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t todoDo) Limit(limit int) *todoDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t todoDo) Offset(offset int) *todoDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t todoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *todoDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t todoDo) Unscoped() *todoDo {
	return t.withDO(t.DO.Unscoped())
}

func (t todoDo) Create(values ...*model.Todo) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t todoDo) CreateInBatches(values []*model.Todo, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t todoDo) Save(values ...*model.Todo) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t todoDo) First() (*model.Todo, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Todo), nil
	}
}

func (t todoDo) Take() (*model.Todo, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Todo), nil
	}
}

func (t todoDo) Last() (*model.Todo, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Todo), nil
	}
}

func (t todoDo) Find() ([]*model.Todo, error) {
	result, err := t.DO.Find()
	return result.([]*model.Todo), err
}

func (t todoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Todo, err error) {
	buf := make([]*model.Todo, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t todoDo) FindInBatches(result *[]*model.Todo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t todoDo) Attrs(attrs ...field.AssignExpr) *todoDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t todoDo) Assign(attrs ...field.AssignExpr) *todoDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t todoDo) Joins(fields ...field.RelationField) *todoDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t todoDo) Preload(fields ...field.RelationField) *todoDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t todoDo) FirstOrInit() (*model.Todo, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Todo), nil
	}
}

func (t todoDo) FirstOrCreate() (*model.Todo, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Todo), nil
	}
}

func (t todoDo) FindByPage(offset int, limit int) (result []*model.Todo, count int64, err error) {
	result, err = t.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = t.Offset(-1).Limit(-1).Count()
	return
}

func (t todoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t todoDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t todoDo) Delete(models ...*model.Todo) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *todoDo) withDO(do gen.Dao) *todoDo {
	t.DO = *do.(*gen.DO)
	return t
}
