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

func newBehavior(db *gorm.DB, opts ...gen.DOOption) behavior {
	_behavior := behavior{}

	_behavior.behaviorDo.UseDB(db, opts...)
	_behavior.behaviorDo.UseModel(&model.Behavior{})

	tableName := _behavior.behaviorDo.TableName()
	_behavior.ALL = field.NewAsterisk(tableName)
	_behavior.ID = field.NewInt32(tableName, "id")
	_behavior.UID = field.NewString(tableName, "uid")
	_behavior.Flag = field.NewString(tableName, "flag")
	_behavior.Count_ = field.NewInt32(tableName, "count")
	_behavior.Extra = field.NewField(tableName, "extra")
	_behavior.CreatedAt = field.NewTime(tableName, "created_at")
	_behavior.UpdatedAt = field.NewTime(tableName, "updated_at")

	_behavior.fillFieldMap()

	return _behavior
}

type behavior struct {
	behaviorDo

	ALL       field.Asterisk
	ID        field.Int32
	UID       field.String
	Flag      field.String
	Count_    field.Int32
	Extra     field.Field
	CreatedAt field.Time
	UpdatedAt field.Time

	fieldMap map[string]field.Expr
}

func (b behavior) Table(newTableName string) *behavior {
	b.behaviorDo.UseTable(newTableName)
	return b.updateTableName(newTableName)
}

func (b behavior) As(alias string) *behavior {
	b.behaviorDo.DO = *(b.behaviorDo.As(alias).(*gen.DO))
	return b.updateTableName(alias)
}

func (b *behavior) updateTableName(table string) *behavior {
	b.ALL = field.NewAsterisk(table)
	b.ID = field.NewInt32(table, "id")
	b.UID = field.NewString(table, "uid")
	b.Flag = field.NewString(table, "flag")
	b.Count_ = field.NewInt32(table, "count")
	b.Extra = field.NewField(table, "extra")
	b.CreatedAt = field.NewTime(table, "created_at")
	b.UpdatedAt = field.NewTime(table, "updated_at")

	b.fillFieldMap()

	return b
}

func (b *behavior) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := b.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (b *behavior) fillFieldMap() {
	b.fieldMap = make(map[string]field.Expr, 7)
	b.fieldMap["id"] = b.ID
	b.fieldMap["uid"] = b.UID
	b.fieldMap["flag"] = b.Flag
	b.fieldMap["count"] = b.Count_
	b.fieldMap["extra"] = b.Extra
	b.fieldMap["created_at"] = b.CreatedAt
	b.fieldMap["updated_at"] = b.UpdatedAt
}

func (b behavior) clone(db *gorm.DB) behavior {
	b.behaviorDo.ReplaceConnPool(db.Statement.ConnPool)
	return b
}

func (b behavior) replaceDB(db *gorm.DB) behavior {
	b.behaviorDo.ReplaceDB(db)
	return b
}

type behaviorDo struct{ gen.DO }

func (b behaviorDo) Debug() *behaviorDo {
	return b.withDO(b.DO.Debug())
}

func (b behaviorDo) WithContext(ctx context.Context) *behaviorDo {
	return b.withDO(b.DO.WithContext(ctx))
}

func (b behaviorDo) ReadDB() *behaviorDo {
	return b.Clauses(dbresolver.Read)
}

func (b behaviorDo) WriteDB() *behaviorDo {
	return b.Clauses(dbresolver.Write)
}

func (b behaviorDo) Session(config *gorm.Session) *behaviorDo {
	return b.withDO(b.DO.Session(config))
}

func (b behaviorDo) Clauses(conds ...clause.Expression) *behaviorDo {
	return b.withDO(b.DO.Clauses(conds...))
}

func (b behaviorDo) Returning(value interface{}, columns ...string) *behaviorDo {
	return b.withDO(b.DO.Returning(value, columns...))
}

func (b behaviorDo) Not(conds ...gen.Condition) *behaviorDo {
	return b.withDO(b.DO.Not(conds...))
}

func (b behaviorDo) Or(conds ...gen.Condition) *behaviorDo {
	return b.withDO(b.DO.Or(conds...))
}

func (b behaviorDo) Select(conds ...field.Expr) *behaviorDo {
	return b.withDO(b.DO.Select(conds...))
}

func (b behaviorDo) Where(conds ...gen.Condition) *behaviorDo {
	return b.withDO(b.DO.Where(conds...))
}

func (b behaviorDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *behaviorDo {
	return b.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (b behaviorDo) Order(conds ...field.Expr) *behaviorDo {
	return b.withDO(b.DO.Order(conds...))
}

func (b behaviorDo) Distinct(cols ...field.Expr) *behaviorDo {
	return b.withDO(b.DO.Distinct(cols...))
}

func (b behaviorDo) Omit(cols ...field.Expr) *behaviorDo {
	return b.withDO(b.DO.Omit(cols...))
}

func (b behaviorDo) Join(table schema.Tabler, on ...field.Expr) *behaviorDo {
	return b.withDO(b.DO.Join(table, on...))
}

func (b behaviorDo) LeftJoin(table schema.Tabler, on ...field.Expr) *behaviorDo {
	return b.withDO(b.DO.LeftJoin(table, on...))
}

func (b behaviorDo) RightJoin(table schema.Tabler, on ...field.Expr) *behaviorDo {
	return b.withDO(b.DO.RightJoin(table, on...))
}

func (b behaviorDo) Group(cols ...field.Expr) *behaviorDo {
	return b.withDO(b.DO.Group(cols...))
}

func (b behaviorDo) Having(conds ...gen.Condition) *behaviorDo {
	return b.withDO(b.DO.Having(conds...))
}

func (b behaviorDo) Limit(limit int) *behaviorDo {
	return b.withDO(b.DO.Limit(limit))
}

func (b behaviorDo) Offset(offset int) *behaviorDo {
	return b.withDO(b.DO.Offset(offset))
}

func (b behaviorDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *behaviorDo {
	return b.withDO(b.DO.Scopes(funcs...))
}

func (b behaviorDo) Unscoped() *behaviorDo {
	return b.withDO(b.DO.Unscoped())
}

func (b behaviorDo) Create(values ...*model.Behavior) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Create(values)
}

func (b behaviorDo) CreateInBatches(values []*model.Behavior, batchSize int) error {
	return b.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (b behaviorDo) Save(values ...*model.Behavior) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Save(values)
}

func (b behaviorDo) First() (*model.Behavior, error) {
	if result, err := b.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Behavior), nil
	}
}

func (b behaviorDo) Take() (*model.Behavior, error) {
	if result, err := b.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Behavior), nil
	}
}

func (b behaviorDo) Last() (*model.Behavior, error) {
	if result, err := b.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Behavior), nil
	}
}

func (b behaviorDo) Find() ([]*model.Behavior, error) {
	result, err := b.DO.Find()
	return result.([]*model.Behavior), err
}

func (b behaviorDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Behavior, err error) {
	buf := make([]*model.Behavior, 0, batchSize)
	err = b.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (b behaviorDo) FindInBatches(result *[]*model.Behavior, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return b.DO.FindInBatches(result, batchSize, fc)
}

func (b behaviorDo) Attrs(attrs ...field.AssignExpr) *behaviorDo {
	return b.withDO(b.DO.Attrs(attrs...))
}

func (b behaviorDo) Assign(attrs ...field.AssignExpr) *behaviorDo {
	return b.withDO(b.DO.Assign(attrs...))
}

func (b behaviorDo) Joins(fields ...field.RelationField) *behaviorDo {
	for _, _f := range fields {
		b = *b.withDO(b.DO.Joins(_f))
	}
	return &b
}

func (b behaviorDo) Preload(fields ...field.RelationField) *behaviorDo {
	for _, _f := range fields {
		b = *b.withDO(b.DO.Preload(_f))
	}
	return &b
}

func (b behaviorDo) FirstOrInit() (*model.Behavior, error) {
	if result, err := b.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Behavior), nil
	}
}

func (b behaviorDo) FirstOrCreate() (*model.Behavior, error) {
	if result, err := b.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Behavior), nil
	}
}

func (b behaviorDo) FindByPage(offset int, limit int) (result []*model.Behavior, count int64, err error) {
	result, err = b.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = b.Offset(-1).Limit(-1).Count()
	return
}

func (b behaviorDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = b.Count()
	if err != nil {
		return
	}

	err = b.Offset(offset).Limit(limit).Scan(result)
	return
}

func (b behaviorDo) Scan(result interface{}) (err error) {
	return b.DO.Scan(result)
}

func (b behaviorDo) Delete(models ...*model.Behavior) (result gen.ResultInfo, err error) {
	return b.DO.Delete(models)
}

func (b *behaviorDo) withDO(do gen.Dao) *behaviorDo {
	b.DO = *do.(*gen.DO)
	return b
}
