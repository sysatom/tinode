// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/tinode/chat/server/extra/store/model"
)

func newKeyResult(db *gorm.DB, opts ...gen.DOOption) keyResult {
	_keyResult := keyResult{}

	_keyResult.keyResultDo.UseDB(db, opts...)
	_keyResult.keyResultDo.UseModel(&model.KeyResult{})

	tableName := _keyResult.keyResultDo.TableName()
	_keyResult.ALL = field.NewAsterisk(tableName)
	_keyResult.ID = field.NewInt32(tableName, "id")
	_keyResult.UID = field.NewString(tableName, "uid")
	_keyResult.Topic = field.NewString(tableName, "topic")
	_keyResult.ObjectiveID = field.NewInt32(tableName, "objective_id")
	_keyResult.Sequence = field.NewInt32(tableName, "sequence")
	_keyResult.Title = field.NewString(tableName, "title")
	_keyResult.Memo = field.NewString(tableName, "memo")
	_keyResult.InitialValue = field.NewInt32(tableName, "initial_value")
	_keyResult.TargetValue = field.NewInt32(tableName, "target_value")
	_keyResult.CurrentValue = field.NewInt32(tableName, "current_value")
	_keyResult.ValueMode = field.NewField(tableName, "value_mode")
	_keyResult.Tag = field.NewString(tableName, "tag")
	_keyResult.CreatedAt = field.NewTime(tableName, "created_at")
	_keyResult.UpdatedAt = field.NewTime(tableName, "updated_at")
	_keyResult.KeyResultValues = keyResultHasManyKeyResultValues{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("KeyResultValues", "model.KeyResultValue"),
	}

	_keyResult.fillFieldMap()

	return _keyResult
}

type keyResult struct {
	keyResultDo

	ALL             field.Asterisk
	ID              field.Int32
	UID             field.String
	Topic           field.String
	ObjectiveID     field.Int32
	Sequence        field.Int32
	Title           field.String
	Memo            field.String
	InitialValue    field.Int32
	TargetValue     field.Int32
	CurrentValue    field.Int32
	ValueMode       field.Field
	Tag             field.String
	CreatedAt       field.Time
	UpdatedAt       field.Time
	KeyResultValues keyResultHasManyKeyResultValues

	fieldMap map[string]field.Expr
}

func (k keyResult) Table(newTableName string) *keyResult {
	k.keyResultDo.UseTable(newTableName)
	return k.updateTableName(newTableName)
}

func (k keyResult) As(alias string) *keyResult {
	k.keyResultDo.DO = *(k.keyResultDo.As(alias).(*gen.DO))
	return k.updateTableName(alias)
}

func (k *keyResult) updateTableName(table string) *keyResult {
	k.ALL = field.NewAsterisk(table)
	k.ID = field.NewInt32(table, "id")
	k.UID = field.NewString(table, "uid")
	k.Topic = field.NewString(table, "topic")
	k.ObjectiveID = field.NewInt32(table, "objective_id")
	k.Sequence = field.NewInt32(table, "sequence")
	k.Title = field.NewString(table, "title")
	k.Memo = field.NewString(table, "memo")
	k.InitialValue = field.NewInt32(table, "initial_value")
	k.TargetValue = field.NewInt32(table, "target_value")
	k.CurrentValue = field.NewInt32(table, "current_value")
	k.ValueMode = field.NewField(table, "value_mode")
	k.Tag = field.NewString(table, "tag")
	k.CreatedAt = field.NewTime(table, "created_at")
	k.UpdatedAt = field.NewTime(table, "updated_at")

	k.fillFieldMap()

	return k
}

func (k *keyResult) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := k.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (k *keyResult) fillFieldMap() {
	k.fieldMap = make(map[string]field.Expr, 15)
	k.fieldMap["id"] = k.ID
	k.fieldMap["uid"] = k.UID
	k.fieldMap["topic"] = k.Topic
	k.fieldMap["objective_id"] = k.ObjectiveID
	k.fieldMap["sequence"] = k.Sequence
	k.fieldMap["title"] = k.Title
	k.fieldMap["memo"] = k.Memo
	k.fieldMap["initial_value"] = k.InitialValue
	k.fieldMap["target_value"] = k.TargetValue
	k.fieldMap["current_value"] = k.CurrentValue
	k.fieldMap["value_mode"] = k.ValueMode
	k.fieldMap["tag"] = k.Tag
	k.fieldMap["created_at"] = k.CreatedAt
	k.fieldMap["updated_at"] = k.UpdatedAt

}

func (k keyResult) clone(db *gorm.DB) keyResult {
	k.keyResultDo.ReplaceConnPool(db.Statement.ConnPool)
	return k
}

func (k keyResult) replaceDB(db *gorm.DB) keyResult {
	k.keyResultDo.ReplaceDB(db)
	return k
}

type keyResultHasManyKeyResultValues struct {
	db *gorm.DB

	field.RelationField
}

func (a keyResultHasManyKeyResultValues) Where(conds ...field.Expr) *keyResultHasManyKeyResultValues {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a keyResultHasManyKeyResultValues) WithContext(ctx context.Context) *keyResultHasManyKeyResultValues {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a keyResultHasManyKeyResultValues) Session(session *gorm.Session) *keyResultHasManyKeyResultValues {
	a.db = a.db.Session(session)
	return &a
}

func (a keyResultHasManyKeyResultValues) Model(m *model.KeyResult) *keyResultHasManyKeyResultValuesTx {
	return &keyResultHasManyKeyResultValuesTx{a.db.Model(m).Association(a.Name())}
}

type keyResultHasManyKeyResultValuesTx struct{ tx *gorm.Association }

func (a keyResultHasManyKeyResultValuesTx) Find() (result []*model.KeyResultValue, err error) {
	return result, a.tx.Find(&result)
}

func (a keyResultHasManyKeyResultValuesTx) Append(values ...*model.KeyResultValue) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a keyResultHasManyKeyResultValuesTx) Replace(values ...*model.KeyResultValue) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a keyResultHasManyKeyResultValuesTx) Delete(values ...*model.KeyResultValue) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a keyResultHasManyKeyResultValuesTx) Clear() error {
	return a.tx.Clear()
}

func (a keyResultHasManyKeyResultValuesTx) Count() int64 {
	return a.tx.Count()
}

type keyResultDo struct{ gen.DO }

// GetByID
//
// SELECT * FROM @@table WHERE id=@id
func (k keyResultDo) GetByID(id int) (result model.KeyResult, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM chatbot_key_results WHERE id=? ")

	var executeSQL *gorm.DB
	executeSQL = k.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// TakeByUidAndTopic query data by uid and topic and return one
//
// where("uid=@uid AND topic=@topic")
func (k keyResultDo) TakeByUidAndTopic(uid string, topic string) (result model.KeyResult, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, uid)
	params = append(params, topic)
	generateSQL.WriteString("uid=? AND topic=? ")

	var executeSQL *gorm.DB
	executeSQL = k.UnderlyingDB().Where(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// FindByUidAndTopic query data by uid and topic and return array
//
// where("uid=@uid AND topic=@topic")
func (k keyResultDo) FindByUidAndTopic(uid string, topic string) (result []*model.KeyResult, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, uid)
	params = append(params, topic)
	generateSQL.WriteString("uid=? AND topic=? ")

	var executeSQL *gorm.DB
	executeSQL = k.UnderlyingDB().Where(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (k keyResultDo) Debug() *keyResultDo {
	return k.withDO(k.DO.Debug())
}

func (k keyResultDo) WithContext(ctx context.Context) *keyResultDo {
	return k.withDO(k.DO.WithContext(ctx))
}

func (k keyResultDo) ReadDB() *keyResultDo {
	return k.Clauses(dbresolver.Read)
}

func (k keyResultDo) WriteDB() *keyResultDo {
	return k.Clauses(dbresolver.Write)
}

func (k keyResultDo) Session(config *gorm.Session) *keyResultDo {
	return k.withDO(k.DO.Session(config))
}

func (k keyResultDo) Clauses(conds ...clause.Expression) *keyResultDo {
	return k.withDO(k.DO.Clauses(conds...))
}

func (k keyResultDo) Returning(value interface{}, columns ...string) *keyResultDo {
	return k.withDO(k.DO.Returning(value, columns...))
}

func (k keyResultDo) Not(conds ...gen.Condition) *keyResultDo {
	return k.withDO(k.DO.Not(conds...))
}

func (k keyResultDo) Or(conds ...gen.Condition) *keyResultDo {
	return k.withDO(k.DO.Or(conds...))
}

func (k keyResultDo) Select(conds ...field.Expr) *keyResultDo {
	return k.withDO(k.DO.Select(conds...))
}

func (k keyResultDo) Where(conds ...gen.Condition) *keyResultDo {
	return k.withDO(k.DO.Where(conds...))
}

func (k keyResultDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *keyResultDo {
	return k.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (k keyResultDo) Order(conds ...field.Expr) *keyResultDo {
	return k.withDO(k.DO.Order(conds...))
}

func (k keyResultDo) Distinct(cols ...field.Expr) *keyResultDo {
	return k.withDO(k.DO.Distinct(cols...))
}

func (k keyResultDo) Omit(cols ...field.Expr) *keyResultDo {
	return k.withDO(k.DO.Omit(cols...))
}

func (k keyResultDo) Join(table schema.Tabler, on ...field.Expr) *keyResultDo {
	return k.withDO(k.DO.Join(table, on...))
}

func (k keyResultDo) LeftJoin(table schema.Tabler, on ...field.Expr) *keyResultDo {
	return k.withDO(k.DO.LeftJoin(table, on...))
}

func (k keyResultDo) RightJoin(table schema.Tabler, on ...field.Expr) *keyResultDo {
	return k.withDO(k.DO.RightJoin(table, on...))
}

func (k keyResultDo) Group(cols ...field.Expr) *keyResultDo {
	return k.withDO(k.DO.Group(cols...))
}

func (k keyResultDo) Having(conds ...gen.Condition) *keyResultDo {
	return k.withDO(k.DO.Having(conds...))
}

func (k keyResultDo) Limit(limit int) *keyResultDo {
	return k.withDO(k.DO.Limit(limit))
}

func (k keyResultDo) Offset(offset int) *keyResultDo {
	return k.withDO(k.DO.Offset(offset))
}

func (k keyResultDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *keyResultDo {
	return k.withDO(k.DO.Scopes(funcs...))
}

func (k keyResultDo) Unscoped() *keyResultDo {
	return k.withDO(k.DO.Unscoped())
}

func (k keyResultDo) Create(values ...*model.KeyResult) error {
	if len(values) == 0 {
		return nil
	}
	return k.DO.Create(values)
}

func (k keyResultDo) CreateInBatches(values []*model.KeyResult, batchSize int) error {
	return k.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (k keyResultDo) Save(values ...*model.KeyResult) error {
	if len(values) == 0 {
		return nil
	}
	return k.DO.Save(values)
}

func (k keyResultDo) First() (*model.KeyResult, error) {
	if result, err := k.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeyResult), nil
	}
}

func (k keyResultDo) Take() (*model.KeyResult, error) {
	if result, err := k.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeyResult), nil
	}
}

func (k keyResultDo) Last() (*model.KeyResult, error) {
	if result, err := k.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeyResult), nil
	}
}

func (k keyResultDo) Find() ([]*model.KeyResult, error) {
	result, err := k.DO.Find()
	return result.([]*model.KeyResult), err
}

func (k keyResultDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.KeyResult, err error) {
	buf := make([]*model.KeyResult, 0, batchSize)
	err = k.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (k keyResultDo) FindInBatches(result *[]*model.KeyResult, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return k.DO.FindInBatches(result, batchSize, fc)
}

func (k keyResultDo) Attrs(attrs ...field.AssignExpr) *keyResultDo {
	return k.withDO(k.DO.Attrs(attrs...))
}

func (k keyResultDo) Assign(attrs ...field.AssignExpr) *keyResultDo {
	return k.withDO(k.DO.Assign(attrs...))
}

func (k keyResultDo) Joins(fields ...field.RelationField) *keyResultDo {
	for _, _f := range fields {
		k = *k.withDO(k.DO.Joins(_f))
	}
	return &k
}

func (k keyResultDo) Preload(fields ...field.RelationField) *keyResultDo {
	for _, _f := range fields {
		k = *k.withDO(k.DO.Preload(_f))
	}
	return &k
}

func (k keyResultDo) FirstOrInit() (*model.KeyResult, error) {
	if result, err := k.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeyResult), nil
	}
}

func (k keyResultDo) FirstOrCreate() (*model.KeyResult, error) {
	if result, err := k.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.KeyResult), nil
	}
}

func (k keyResultDo) FindByPage(offset int, limit int) (result []*model.KeyResult, count int64, err error) {
	result, err = k.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = k.Offset(-1).Limit(-1).Count()
	return
}

func (k keyResultDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = k.Count()
	if err != nil {
		return
	}

	err = k.Offset(offset).Limit(limit).Scan(result)
	return
}

func (k keyResultDo) Scan(result interface{}) (err error) {
	return k.DO.Scan(result)
}

func (k keyResultDo) Delete(models ...*model.KeyResult) (result gen.ResultInfo, err error) {
	return k.DO.Delete(models)
}

func (k *keyResultDo) withDO(do gen.Dao) *keyResultDo {
	k.DO = *do.(*gen.DO)
	return k
}
