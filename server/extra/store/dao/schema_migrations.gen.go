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

func newSchemaMigration(db *gorm.DB, opts ...gen.DOOption) schemaMigration {
	_schemaMigration := schemaMigration{}

	_schemaMigration.schemaMigrationDo.UseDB(db, opts...)
	_schemaMigration.schemaMigrationDo.UseModel(&model.SchemaMigration{})

	tableName := _schemaMigration.schemaMigrationDo.TableName()
	_schemaMigration.ALL = field.NewAsterisk(tableName)
	_schemaMigration.Version = field.NewInt64(tableName, "version")
	_schemaMigration.Dirty = field.NewBool(tableName, "dirty")

	_schemaMigration.fillFieldMap()

	return _schemaMigration
}

type schemaMigration struct {
	schemaMigrationDo

	ALL     field.Asterisk
	Version field.Int64
	Dirty   field.Bool

	fieldMap map[string]field.Expr
}

func (s schemaMigration) Table(newTableName string) *schemaMigration {
	s.schemaMigrationDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s schemaMigration) As(alias string) *schemaMigration {
	s.schemaMigrationDo.DO = *(s.schemaMigrationDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *schemaMigration) updateTableName(table string) *schemaMigration {
	s.ALL = field.NewAsterisk(table)
	s.Version = field.NewInt64(table, "version")
	s.Dirty = field.NewBool(table, "dirty")

	s.fillFieldMap()

	return s
}

func (s *schemaMigration) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *schemaMigration) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 2)
	s.fieldMap["version"] = s.Version
	s.fieldMap["dirty"] = s.Dirty
}

func (s schemaMigration) clone(db *gorm.DB) schemaMigration {
	s.schemaMigrationDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s schemaMigration) replaceDB(db *gorm.DB) schemaMigration {
	s.schemaMigrationDo.ReplaceDB(db)
	return s
}

type schemaMigrationDo struct{ gen.DO }

func (s schemaMigrationDo) Debug() *schemaMigrationDo {
	return s.withDO(s.DO.Debug())
}

func (s schemaMigrationDo) WithContext(ctx context.Context) *schemaMigrationDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s schemaMigrationDo) ReadDB() *schemaMigrationDo {
	return s.Clauses(dbresolver.Read)
}

func (s schemaMigrationDo) WriteDB() *schemaMigrationDo {
	return s.Clauses(dbresolver.Write)
}

func (s schemaMigrationDo) Session(config *gorm.Session) *schemaMigrationDo {
	return s.withDO(s.DO.Session(config))
}

func (s schemaMigrationDo) Clauses(conds ...clause.Expression) *schemaMigrationDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s schemaMigrationDo) Returning(value interface{}, columns ...string) *schemaMigrationDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s schemaMigrationDo) Not(conds ...gen.Condition) *schemaMigrationDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s schemaMigrationDo) Or(conds ...gen.Condition) *schemaMigrationDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s schemaMigrationDo) Select(conds ...field.Expr) *schemaMigrationDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s schemaMigrationDo) Where(conds ...gen.Condition) *schemaMigrationDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s schemaMigrationDo) Order(conds ...field.Expr) *schemaMigrationDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s schemaMigrationDo) Distinct(cols ...field.Expr) *schemaMigrationDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s schemaMigrationDo) Omit(cols ...field.Expr) *schemaMigrationDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s schemaMigrationDo) Join(table schema.Tabler, on ...field.Expr) *schemaMigrationDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s schemaMigrationDo) LeftJoin(table schema.Tabler, on ...field.Expr) *schemaMigrationDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s schemaMigrationDo) RightJoin(table schema.Tabler, on ...field.Expr) *schemaMigrationDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s schemaMigrationDo) Group(cols ...field.Expr) *schemaMigrationDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s schemaMigrationDo) Having(conds ...gen.Condition) *schemaMigrationDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s schemaMigrationDo) Limit(limit int) *schemaMigrationDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s schemaMigrationDo) Offset(offset int) *schemaMigrationDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s schemaMigrationDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *schemaMigrationDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s schemaMigrationDo) Unscoped() *schemaMigrationDo {
	return s.withDO(s.DO.Unscoped())
}

func (s schemaMigrationDo) Create(values ...*model.SchemaMigration) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s schemaMigrationDo) CreateInBatches(values []*model.SchemaMigration, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s schemaMigrationDo) Save(values ...*model.SchemaMigration) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s schemaMigrationDo) First() (*model.SchemaMigration, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.SchemaMigration), nil
	}
}

func (s schemaMigrationDo) Take() (*model.SchemaMigration, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.SchemaMigration), nil
	}
}

func (s schemaMigrationDo) Last() (*model.SchemaMigration, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.SchemaMigration), nil
	}
}

func (s schemaMigrationDo) Find() ([]*model.SchemaMigration, error) {
	result, err := s.DO.Find()
	return result.([]*model.SchemaMigration), err
}

func (s schemaMigrationDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.SchemaMigration, err error) {
	buf := make([]*model.SchemaMigration, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s schemaMigrationDo) FindInBatches(result *[]*model.SchemaMigration, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s schemaMigrationDo) Attrs(attrs ...field.AssignExpr) *schemaMigrationDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s schemaMigrationDo) Assign(attrs ...field.AssignExpr) *schemaMigrationDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s schemaMigrationDo) Joins(fields ...field.RelationField) *schemaMigrationDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s schemaMigrationDo) Preload(fields ...field.RelationField) *schemaMigrationDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s schemaMigrationDo) FirstOrInit() (*model.SchemaMigration, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.SchemaMigration), nil
	}
}

func (s schemaMigrationDo) FirstOrCreate() (*model.SchemaMigration, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.SchemaMigration), nil
	}
}

func (s schemaMigrationDo) FindByPage(offset int, limit int) (result []*model.SchemaMigration, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s schemaMigrationDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s schemaMigrationDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s schemaMigrationDo) Delete(models ...*model.SchemaMigration) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *schemaMigrationDo) withDO(do gen.Dao) *schemaMigrationDo {
	s.DO = *do.(*gen.DO)
	return s
}
