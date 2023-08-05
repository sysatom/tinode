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

func newSubscription(db *gorm.DB, opts ...gen.DOOption) subscription {
	_subscription := subscription{}

	_subscription.subscriptionDo.UseDB(db, opts...)
	_subscription.subscriptionDo.UseModel(&model.Subscription{})

	tableName := _subscription.subscriptionDo.TableName()
	_subscription.ALL = field.NewAsterisk(tableName)
	_subscription.ID = field.NewInt32(tableName, "id")
	_subscription.Createdat = field.NewTime(tableName, "createdat")
	_subscription.Updatedat = field.NewTime(tableName, "updatedat")
	_subscription.Deletedat = field.NewTime(tableName, "deletedat")
	_subscription.Userid = field.NewInt64(tableName, "userid")
	_subscription.Topic = field.NewString(tableName, "topic")
	_subscription.Delid = field.NewInt32(tableName, "delid")
	_subscription.Recvseqid = field.NewInt32(tableName, "recvseqid")
	_subscription.Readseqid = field.NewInt32(tableName, "readseqid")
	_subscription.Modewant = field.NewString(tableName, "modewant")
	_subscription.Modegiven = field.NewString(tableName, "modegiven")
	_subscription.Private = field.NewString(tableName, "private")

	_subscription.fillFieldMap()

	return _subscription
}

type subscription struct {
	subscriptionDo

	ALL       field.Asterisk
	ID        field.Int32
	Createdat field.Time
	Updatedat field.Time
	Deletedat field.Time
	Userid    field.Int64
	Topic     field.String
	Delid     field.Int32
	Recvseqid field.Int32
	Readseqid field.Int32
	Modewant  field.String
	Modegiven field.String
	Private   field.String

	fieldMap map[string]field.Expr
}

func (s subscription) Table(newTableName string) *subscription {
	s.subscriptionDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s subscription) As(alias string) *subscription {
	s.subscriptionDo.DO = *(s.subscriptionDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *subscription) updateTableName(table string) *subscription {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewInt32(table, "id")
	s.Createdat = field.NewTime(table, "createdat")
	s.Updatedat = field.NewTime(table, "updatedat")
	s.Deletedat = field.NewTime(table, "deletedat")
	s.Userid = field.NewInt64(table, "userid")
	s.Topic = field.NewString(table, "topic")
	s.Delid = field.NewInt32(table, "delid")
	s.Recvseqid = field.NewInt32(table, "recvseqid")
	s.Readseqid = field.NewInt32(table, "readseqid")
	s.Modewant = field.NewString(table, "modewant")
	s.Modegiven = field.NewString(table, "modegiven")
	s.Private = field.NewString(table, "private")

	s.fillFieldMap()

	return s
}

func (s *subscription) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *subscription) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 12)
	s.fieldMap["id"] = s.ID
	s.fieldMap["createdat"] = s.Createdat
	s.fieldMap["updatedat"] = s.Updatedat
	s.fieldMap["deletedat"] = s.Deletedat
	s.fieldMap["userid"] = s.Userid
	s.fieldMap["topic"] = s.Topic
	s.fieldMap["delid"] = s.Delid
	s.fieldMap["recvseqid"] = s.Recvseqid
	s.fieldMap["readseqid"] = s.Readseqid
	s.fieldMap["modewant"] = s.Modewant
	s.fieldMap["modegiven"] = s.Modegiven
	s.fieldMap["private"] = s.Private
}

func (s subscription) clone(db *gorm.DB) subscription {
	s.subscriptionDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s subscription) replaceDB(db *gorm.DB) subscription {
	s.subscriptionDo.ReplaceDB(db)
	return s
}

type subscriptionDo struct{ gen.DO }

func (s subscriptionDo) Debug() *subscriptionDo {
	return s.withDO(s.DO.Debug())
}

func (s subscriptionDo) WithContext(ctx context.Context) *subscriptionDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s subscriptionDo) ReadDB() *subscriptionDo {
	return s.Clauses(dbresolver.Read)
}

func (s subscriptionDo) WriteDB() *subscriptionDo {
	return s.Clauses(dbresolver.Write)
}

func (s subscriptionDo) Session(config *gorm.Session) *subscriptionDo {
	return s.withDO(s.DO.Session(config))
}

func (s subscriptionDo) Clauses(conds ...clause.Expression) *subscriptionDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s subscriptionDo) Returning(value interface{}, columns ...string) *subscriptionDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s subscriptionDo) Not(conds ...gen.Condition) *subscriptionDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s subscriptionDo) Or(conds ...gen.Condition) *subscriptionDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s subscriptionDo) Select(conds ...field.Expr) *subscriptionDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s subscriptionDo) Where(conds ...gen.Condition) *subscriptionDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s subscriptionDo) Order(conds ...field.Expr) *subscriptionDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s subscriptionDo) Distinct(cols ...field.Expr) *subscriptionDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s subscriptionDo) Omit(cols ...field.Expr) *subscriptionDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s subscriptionDo) Join(table schema.Tabler, on ...field.Expr) *subscriptionDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s subscriptionDo) LeftJoin(table schema.Tabler, on ...field.Expr) *subscriptionDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s subscriptionDo) RightJoin(table schema.Tabler, on ...field.Expr) *subscriptionDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s subscriptionDo) Group(cols ...field.Expr) *subscriptionDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s subscriptionDo) Having(conds ...gen.Condition) *subscriptionDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s subscriptionDo) Limit(limit int) *subscriptionDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s subscriptionDo) Offset(offset int) *subscriptionDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s subscriptionDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *subscriptionDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s subscriptionDo) Unscoped() *subscriptionDo {
	return s.withDO(s.DO.Unscoped())
}

func (s subscriptionDo) Create(values ...*model.Subscription) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s subscriptionDo) CreateInBatches(values []*model.Subscription, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s subscriptionDo) Save(values ...*model.Subscription) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s subscriptionDo) First() (*model.Subscription, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Subscription), nil
	}
}

func (s subscriptionDo) Take() (*model.Subscription, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Subscription), nil
	}
}

func (s subscriptionDo) Last() (*model.Subscription, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Subscription), nil
	}
}

func (s subscriptionDo) Find() ([]*model.Subscription, error) {
	result, err := s.DO.Find()
	return result.([]*model.Subscription), err
}

func (s subscriptionDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Subscription, err error) {
	buf := make([]*model.Subscription, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s subscriptionDo) FindInBatches(result *[]*model.Subscription, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s subscriptionDo) Attrs(attrs ...field.AssignExpr) *subscriptionDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s subscriptionDo) Assign(attrs ...field.AssignExpr) *subscriptionDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s subscriptionDo) Joins(fields ...field.RelationField) *subscriptionDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s subscriptionDo) Preload(fields ...field.RelationField) *subscriptionDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s subscriptionDo) FirstOrInit() (*model.Subscription, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Subscription), nil
	}
}

func (s subscriptionDo) FirstOrCreate() (*model.Subscription, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Subscription), nil
	}
}

func (s subscriptionDo) FindByPage(offset int, limit int) (result []*model.Subscription, count int64, err error) {
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

func (s subscriptionDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s subscriptionDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s subscriptionDo) Delete(models ...*model.Subscription) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *subscriptionDo) withDO(do gen.Dao) *subscriptionDo {
	s.DO = *do.(*gen.DO)
	return s
}
