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

func newUrl(db *gorm.DB, opts ...gen.DOOption) url {
	_url := url{}

	_url.urlDo.UseDB(db, opts...)
	_url.urlDo.UseModel(&model.Url{})

	tableName := _url.urlDo.TableName()
	_url.ALL = field.NewAsterisk(tableName)
	_url.ID = field.NewUint64(tableName, "id")
	_url.Flag = field.NewString(tableName, "flag")
	_url.Url = field.NewString(tableName, "url")
	_url.State = field.NewInt(tableName, "state")
	_url.ViewCount = field.NewInt(tableName, "view_count")
	_url.CreatedAt = field.NewTime(tableName, "created_at")
	_url.UpdatedAt = field.NewTime(tableName, "updated_at")

	_url.fillFieldMap()

	return _url
}

type url struct {
	urlDo

	ALL       field.Asterisk
	ID        field.Uint64
	Flag      field.String
	Url       field.String
	State     field.Int
	ViewCount field.Int
	CreatedAt field.Time
	UpdatedAt field.Time

	fieldMap map[string]field.Expr
}

func (u url) Table(newTableName string) *url {
	u.urlDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u url) As(alias string) *url {
	u.urlDo.DO = *(u.urlDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *url) updateTableName(table string) *url {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewUint64(table, "id")
	u.Flag = field.NewString(table, "flag")
	u.Url = field.NewString(table, "url")
	u.State = field.NewInt(table, "state")
	u.ViewCount = field.NewInt(table, "view_count")
	u.CreatedAt = field.NewTime(table, "created_at")
	u.UpdatedAt = field.NewTime(table, "updated_at")

	u.fillFieldMap()

	return u
}

func (u *url) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *url) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 7)
	u.fieldMap["id"] = u.ID
	u.fieldMap["flag"] = u.Flag
	u.fieldMap["url"] = u.Url
	u.fieldMap["state"] = u.State
	u.fieldMap["view_count"] = u.ViewCount
	u.fieldMap["created_at"] = u.CreatedAt
	u.fieldMap["updated_at"] = u.UpdatedAt
}

func (u url) clone(db *gorm.DB) url {
	u.urlDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u url) replaceDB(db *gorm.DB) url {
	u.urlDo.ReplaceDB(db)
	return u
}

type urlDo struct{ gen.DO }

func (u urlDo) Debug() *urlDo {
	return u.withDO(u.DO.Debug())
}

func (u urlDo) WithContext(ctx context.Context) *urlDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u urlDo) ReadDB() *urlDo {
	return u.Clauses(dbresolver.Read)
}

func (u urlDo) WriteDB() *urlDo {
	return u.Clauses(dbresolver.Write)
}

func (u urlDo) Session(config *gorm.Session) *urlDo {
	return u.withDO(u.DO.Session(config))
}

func (u urlDo) Clauses(conds ...clause.Expression) *urlDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u urlDo) Returning(value interface{}, columns ...string) *urlDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u urlDo) Not(conds ...gen.Condition) *urlDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u urlDo) Or(conds ...gen.Condition) *urlDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u urlDo) Select(conds ...field.Expr) *urlDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u urlDo) Where(conds ...gen.Condition) *urlDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u urlDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *urlDo {
	return u.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (u urlDo) Order(conds ...field.Expr) *urlDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u urlDo) Distinct(cols ...field.Expr) *urlDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u urlDo) Omit(cols ...field.Expr) *urlDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u urlDo) Join(table schema.Tabler, on ...field.Expr) *urlDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u urlDo) LeftJoin(table schema.Tabler, on ...field.Expr) *urlDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u urlDo) RightJoin(table schema.Tabler, on ...field.Expr) *urlDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u urlDo) Group(cols ...field.Expr) *urlDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u urlDo) Having(conds ...gen.Condition) *urlDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u urlDo) Limit(limit int) *urlDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u urlDo) Offset(offset int) *urlDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u urlDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *urlDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u urlDo) Unscoped() *urlDo {
	return u.withDO(u.DO.Unscoped())
}

func (u urlDo) Create(values ...*model.Url) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u urlDo) CreateInBatches(values []*model.Url, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u urlDo) Save(values ...*model.Url) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u urlDo) First() (*model.Url, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Url), nil
	}
}

func (u urlDo) Take() (*model.Url, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Url), nil
	}
}

func (u urlDo) Last() (*model.Url, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Url), nil
	}
}

func (u urlDo) Find() ([]*model.Url, error) {
	result, err := u.DO.Find()
	return result.([]*model.Url), err
}

func (u urlDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Url, err error) {
	buf := make([]*model.Url, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u urlDo) FindInBatches(result *[]*model.Url, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u urlDo) Attrs(attrs ...field.AssignExpr) *urlDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u urlDo) Assign(attrs ...field.AssignExpr) *urlDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u urlDo) Joins(fields ...field.RelationField) *urlDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u urlDo) Preload(fields ...field.RelationField) *urlDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u urlDo) FirstOrInit() (*model.Url, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Url), nil
	}
}

func (u urlDo) FirstOrCreate() (*model.Url, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Url), nil
	}
}

func (u urlDo) FindByPage(offset int, limit int) (result []*model.Url, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u urlDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u urlDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u urlDo) Delete(models ...*model.Url) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *urlDo) withDO(do gen.Dao) *urlDo {
	u.DO = *do.(*gen.DO)
	return u
}
