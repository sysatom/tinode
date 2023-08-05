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

func newDellog(db *gorm.DB, opts ...gen.DOOption) dellog {
	_dellog := dellog{}

	_dellog.dellogDo.UseDB(db, opts...)
	_dellog.dellogDo.UseModel(&model.Dellog{})

	tableName := _dellog.dellogDo.TableName()
	_dellog.ALL = field.NewAsterisk(tableName)
	_dellog.ID = field.NewInt32(tableName, "id")
	_dellog.Topic = field.NewString(tableName, "topic")
	_dellog.Deletedfor = field.NewInt64(tableName, "deletedfor")
	_dellog.Delid = field.NewInt32(tableName, "delid")
	_dellog.Low = field.NewInt32(tableName, "low")
	_dellog.Hi = field.NewInt32(tableName, "hi")

	_dellog.fillFieldMap()

	return _dellog
}

type dellog struct {
	dellogDo

	ALL        field.Asterisk
	ID         field.Int32
	Topic      field.String
	Deletedfor field.Int64
	Delid      field.Int32
	Low        field.Int32
	Hi         field.Int32

	fieldMap map[string]field.Expr
}

func (d dellog) Table(newTableName string) *dellog {
	d.dellogDo.UseTable(newTableName)
	return d.updateTableName(newTableName)
}

func (d dellog) As(alias string) *dellog {
	d.dellogDo.DO = *(d.dellogDo.As(alias).(*gen.DO))
	return d.updateTableName(alias)
}

func (d *dellog) updateTableName(table string) *dellog {
	d.ALL = field.NewAsterisk(table)
	d.ID = field.NewInt32(table, "id")
	d.Topic = field.NewString(table, "topic")
	d.Deletedfor = field.NewInt64(table, "deletedfor")
	d.Delid = field.NewInt32(table, "delid")
	d.Low = field.NewInt32(table, "low")
	d.Hi = field.NewInt32(table, "hi")

	d.fillFieldMap()

	return d
}

func (d *dellog) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := d.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (d *dellog) fillFieldMap() {
	d.fieldMap = make(map[string]field.Expr, 6)
	d.fieldMap["id"] = d.ID
	d.fieldMap["topic"] = d.Topic
	d.fieldMap["deletedfor"] = d.Deletedfor
	d.fieldMap["delid"] = d.Delid
	d.fieldMap["low"] = d.Low
	d.fieldMap["hi"] = d.Hi
}

func (d dellog) clone(db *gorm.DB) dellog {
	d.dellogDo.ReplaceConnPool(db.Statement.ConnPool)
	return d
}

func (d dellog) replaceDB(db *gorm.DB) dellog {
	d.dellogDo.ReplaceDB(db)
	return d
}

type dellogDo struct{ gen.DO }

func (d dellogDo) Debug() *dellogDo {
	return d.withDO(d.DO.Debug())
}

func (d dellogDo) WithContext(ctx context.Context) *dellogDo {
	return d.withDO(d.DO.WithContext(ctx))
}

func (d dellogDo) ReadDB() *dellogDo {
	return d.Clauses(dbresolver.Read)
}

func (d dellogDo) WriteDB() *dellogDo {
	return d.Clauses(dbresolver.Write)
}

func (d dellogDo) Session(config *gorm.Session) *dellogDo {
	return d.withDO(d.DO.Session(config))
}

func (d dellogDo) Clauses(conds ...clause.Expression) *dellogDo {
	return d.withDO(d.DO.Clauses(conds...))
}

func (d dellogDo) Returning(value interface{}, columns ...string) *dellogDo {
	return d.withDO(d.DO.Returning(value, columns...))
}

func (d dellogDo) Not(conds ...gen.Condition) *dellogDo {
	return d.withDO(d.DO.Not(conds...))
}

func (d dellogDo) Or(conds ...gen.Condition) *dellogDo {
	return d.withDO(d.DO.Or(conds...))
}

func (d dellogDo) Select(conds ...field.Expr) *dellogDo {
	return d.withDO(d.DO.Select(conds...))
}

func (d dellogDo) Where(conds ...gen.Condition) *dellogDo {
	return d.withDO(d.DO.Where(conds...))
}

func (d dellogDo) Order(conds ...field.Expr) *dellogDo {
	return d.withDO(d.DO.Order(conds...))
}

func (d dellogDo) Distinct(cols ...field.Expr) *dellogDo {
	return d.withDO(d.DO.Distinct(cols...))
}

func (d dellogDo) Omit(cols ...field.Expr) *dellogDo {
	return d.withDO(d.DO.Omit(cols...))
}

func (d dellogDo) Join(table schema.Tabler, on ...field.Expr) *dellogDo {
	return d.withDO(d.DO.Join(table, on...))
}

func (d dellogDo) LeftJoin(table schema.Tabler, on ...field.Expr) *dellogDo {
	return d.withDO(d.DO.LeftJoin(table, on...))
}

func (d dellogDo) RightJoin(table schema.Tabler, on ...field.Expr) *dellogDo {
	return d.withDO(d.DO.RightJoin(table, on...))
}

func (d dellogDo) Group(cols ...field.Expr) *dellogDo {
	return d.withDO(d.DO.Group(cols...))
}

func (d dellogDo) Having(conds ...gen.Condition) *dellogDo {
	return d.withDO(d.DO.Having(conds...))
}

func (d dellogDo) Limit(limit int) *dellogDo {
	return d.withDO(d.DO.Limit(limit))
}

func (d dellogDo) Offset(offset int) *dellogDo {
	return d.withDO(d.DO.Offset(offset))
}

func (d dellogDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *dellogDo {
	return d.withDO(d.DO.Scopes(funcs...))
}

func (d dellogDo) Unscoped() *dellogDo {
	return d.withDO(d.DO.Unscoped())
}

func (d dellogDo) Create(values ...*model.Dellog) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Create(values)
}

func (d dellogDo) CreateInBatches(values []*model.Dellog, batchSize int) error {
	return d.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (d dellogDo) Save(values ...*model.Dellog) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Save(values)
}

func (d dellogDo) First() (*model.Dellog, error) {
	if result, err := d.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Dellog), nil
	}
}

func (d dellogDo) Take() (*model.Dellog, error) {
	if result, err := d.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Dellog), nil
	}
}

func (d dellogDo) Last() (*model.Dellog, error) {
	if result, err := d.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Dellog), nil
	}
}

func (d dellogDo) Find() ([]*model.Dellog, error) {
	result, err := d.DO.Find()
	return result.([]*model.Dellog), err
}

func (d dellogDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Dellog, err error) {
	buf := make([]*model.Dellog, 0, batchSize)
	err = d.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (d dellogDo) FindInBatches(result *[]*model.Dellog, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return d.DO.FindInBatches(result, batchSize, fc)
}

func (d dellogDo) Attrs(attrs ...field.AssignExpr) *dellogDo {
	return d.withDO(d.DO.Attrs(attrs...))
}

func (d dellogDo) Assign(attrs ...field.AssignExpr) *dellogDo {
	return d.withDO(d.DO.Assign(attrs...))
}

func (d dellogDo) Joins(fields ...field.RelationField) *dellogDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Joins(_f))
	}
	return &d
}

func (d dellogDo) Preload(fields ...field.RelationField) *dellogDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Preload(_f))
	}
	return &d
}

func (d dellogDo) FirstOrInit() (*model.Dellog, error) {
	if result, err := d.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Dellog), nil
	}
}

func (d dellogDo) FirstOrCreate() (*model.Dellog, error) {
	if result, err := d.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Dellog), nil
	}
}

func (d dellogDo) FindByPage(offset int, limit int) (result []*model.Dellog, count int64, err error) {
	result, err = d.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = d.Offset(-1).Limit(-1).Count()
	return
}

func (d dellogDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = d.Count()
	if err != nil {
		return
	}

	err = d.Offset(offset).Limit(limit).Scan(result)
	return
}

func (d dellogDo) Scan(result interface{}) (err error) {
	return d.DO.Scan(result)
}

func (d dellogDo) Delete(models ...*model.Dellog) (result gen.ResultInfo, err error) {
	return d.DO.Delete(models)
}

func (d *dellogDo) withDO(do gen.Dao) *dellogDo {
	d.DO = *do.(*gen.DO)
	return d
}
