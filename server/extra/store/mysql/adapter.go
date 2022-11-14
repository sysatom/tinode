//go:build mysql
// +build mysql

package mysql

import (
	"database/sql"
	"errors"
	"github.com/tinode/chat/server/db/mysql"
	"github.com/tinode/chat/server/extra/locker"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/store/types"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
)

const adapterName = "mysql"

type adapter struct {
	db *gorm.DB
}

func (a *adapter) Open() error {
	db, err := gorm.Open(mysqlDriver.New(mysqlDriver.Config{Conn: mysql.RawDB}), &gorm.Config{})
	if err != nil {
		return err
	}
	a.db = db
	return nil
}

func (a *adapter) IsOpen() bool {
	return a.db != nil
}

func (a *adapter) Close() error {
	rawDB, err := a.db.DB()
	if err != nil {
		return err
	}
	return rawDB.Close()
}

func (a *adapter) GetName() string {
	return adapterName
}

func (a *adapter) Stats() interface{} {
	if a.db == nil {
		return nil
	}
	rawDB, err := a.db.DB()
	if err != nil {
		return err
	}
	return rawDB.Stats()
}

func (a *adapter) DataSet(uid types.Uid, topic, key string, value model.JSON) error {
	var find model.Data
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `key` = ?", uid.UserId(), topic, key).First(&find).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if find.ID > 0 {
		return a.db.
			Model(&model.Config{}).
			Where("`uid` = ? AND `topic` = ? AND `key` = ?", uid.UserId(), topic, key).
			Update("value", value).Error
	} else {
		return a.db.Create(&model.Config{
			Uid:   uid.UserId(),
			Topic: topic,
			Key:   key,
			Value: value,
		}).Error
	}
}

func (a *adapter) DataGet(uid types.Uid, topic, key string) (model.JSON, error) {
	var find model.Data
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `key` = ?", uid.UserId(), topic, key).First(&find).Error
	if err != nil {
		return nil, err
	}
	return find.Value, nil
}

func (a *adapter) ConfigSet(uid types.Uid, topic, key string, value model.JSON) error {
	var find model.Config
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `key` = ?", uid.UserId(), topic, key).First(&find).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if find.ID > 0 {
		return a.db.
			Model(&model.Config{}).
			Where("`uid` = ? AND `topic` = ? AND `key` = ?", uid.UserId(), topic, key).
			Update("value", value).Error
	} else {
		return a.db.Create(&model.Config{
			Uid:   uid.UserId(),
			Topic: topic,
			Key:   key,
			Value: value,
		}).Error
	}
}

func (a *adapter) ConfigGet(uid types.Uid, topic, key string) (model.JSON, error) {
	var find model.Config
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `key` = ?", uid.UserId(), topic, key).First(&find).Error
	if err != nil {
		return nil, err
	}
	return find.Value, nil
}

func (a *adapter) OAuthSet(oauth model.OAuth) error {
	var find model.OAuth
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `type` = ?", oauth.Uid, oauth.Topic, oauth.Type).First(&find).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if find.ID > 0 {
		return a.db.
			Model(&model.OAuth{}).
			Where("`uid` = ? AND `topic` = ? AND `type` = ?", oauth.Uid, oauth.Topic, oauth.Type).
			UpdateColumns(map[string]interface{}{
				"token": oauth.Token,
				"extra": oauth.Extra,
			}).Error
	} else {
		return a.db.Create(&oauth).Error
	}
}

func (a *adapter) OAuthGet(uid types.Uid, topic, t string) (model.OAuth, error) {
	var find model.OAuth
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `type` = ?", uid.UserId(), topic, t).First(&find).Error
	if err != nil {
		return model.OAuth{}, err
	}
	return find, nil
}

func (a *adapter) OAuthGetAvailable(t string) ([]model.OAuth, error) {
	var find []model.OAuth
	err := a.db.Where("`type` = ? AND `token` <> ''", t).Find(&find).Error
	if err != nil {
		return []model.OAuth{}, err
	}
	return find, nil
}

func (a *adapter) FormSet(formId string, form model.Form) error {
	var find model.Form
	err := a.db.Where("`form_id` = ?", formId).First(&find).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if find.ID > 0 {
		return a.db.
			Model(&model.Form{}).
			Where("`form_id` = ?", formId).
			Updates(map[string]interface{}{
				"values": form.Values,
				"state":  form.State,
			}).Error
	} else {
		return a.db.Create(&model.Form{
			FormId: formId,
			Uid:    form.Uid,
			Topic:  form.Topic,
			Schema: form.Schema,
			Values: form.Values,
			State:  form.State,
		}).Error
	}
}

func (a *adapter) FormGet(formId string) (model.Form, error) {
	var find model.Form
	err := a.db.Where("`form_id` = ?", formId).First(&find).Error
	if err != nil {
		return model.Form{}, err
	}
	return find, nil
}

func (a *adapter) PageSet(pageId string, page model.Page) error {
	var find model.Page
	err := a.db.Where("`page_id` = ?", pageId).First(&find).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if find.ID > 0 {
		return a.db.
			Model(&model.Page{}).
			Where("`page_id` = ?", pageId).
			Updates(map[string]interface{}{
				"state": page.State,
			}).Error
	} else {
		return a.db.Create(&page).Error
	}
}

func (a *adapter) PageGet(pageId string) (model.Page, error) {
	var find model.Page
	err := a.db.Where("`page_id` = ?", pageId).First(&find).Error
	if err != nil {
		return model.Page{}, err
	}
	return find, nil
}

func (a *adapter) GetObjectiveByID(id int64) (*model.Objective, error) {
	var objective model.Objective
	err := a.db.Where("id = ?", id).First(&objective).Error
	if err != nil {
		return nil, err
	}
	return &objective, nil
}

func (a *adapter) GetObjectiveBySequence(uid types.Uid, topic string, sequence int64) (*model.Objective, error) {
	var objective model.Objective
	err := a.db.Where("`uid` = ? AND `topic` = ? AND sequence = ?", uid.UserId(), topic, sequence).First(&objective).Error
	if err != nil {
		return nil, err
	}
	return &objective, nil
}

func (a *adapter) ListObjectives(uid types.Uid, topic string) ([]*model.Objective, error) {
	var objectives []*model.Objective
	err := a.db.Where("`uid` = ? AND `topic` = ?", uid.UserId(), topic).Order("id DESC").Find(&objectives).Error
	if err != nil {
		return nil, err
	}
	return objectives, nil
}

func (a *adapter) CreateObjective(objective *model.Objective) (int64, error) {
	locker.Mux.Lock()
	defer locker.Mux.Unlock()

	// sequence
	sequence := int64(0)
	var max model.Objective
	err := a.db.Where("`uid` = ? AND `topic` = ?", objective.Uid, objective.Topic).Order("sequence DESC").Take(&max).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if max.Sequence > 0 {
		sequence = max.Sequence
	}
	sequence += 1

	objective.Sequence = sequence
	err = a.db.Create(&objective).Error
	if err != nil {
		return 0, err
	}
	return objective.Id, nil
}

func (a *adapter) UpdateObjective(objective *model.Objective) error {
	return a.db.Model(&model.Objective{}).
		Where("`uid` = ? AND `topic` = ? AND sequence = ?", objective.Uid, objective.Topic, objective.Sequence).
		UpdateColumns(map[string]interface{}{
			"title":       objective.Title,
			"memo":        objective.Memo,
			"motive":      objective.Motive,
			"feasibility": objective.Feasibility,
			"is_plan":     objective.IsPlan,
			"plan_start":  objective.PlanStart,
			"plan_end":    objective.PlanEnd,
		}).Error
}

func (a *adapter) DeleteObjective(id int64) error {
	return a.db.Where("id = ?", id).Delete(&model.Objective{}).Error
}

func (a *adapter) DeleteObjectiveBySequence(uid types.Uid, topic string, sequence int64) error {
	return a.db.Where("`uid` = ? AND `topic` = ? AND sequence = ?", uid.UserId(), topic, sequence).Delete(&model.Objective{}).Error
}

func (a *adapter) GetKeyResultByID(id int64) (*model.KeyResult, error) {
	var keyResult model.KeyResult
	err := a.db.Where("id = ?", id).First(&keyResult).Error
	if err != nil {
		return nil, err
	}
	return &keyResult, nil
}

func (a *adapter) GetKeyResultBySequence(uid types.Uid, topic string, sequence int64) (*model.KeyResult, error) {
	var keyResult model.KeyResult
	err := a.db.Where("`uid` = ? AND `topic` = ? AND sequence = ?", uid.UserId(), topic, sequence).First(&keyResult).Error
	if err != nil {
		return nil, err
	}
	return &keyResult, nil
}

func (a *adapter) ListKeyResults(uid types.Uid, topic string) ([]*model.KeyResult, error) {
	var keyResult []*model.KeyResult
	err := a.db.Where("`uid` = ? AND `topic` = ?", uid.UserId(), topic).Order("id DESC").Find(&keyResult).Error
	if err != nil {
		return nil, err
	}
	return keyResult, nil
}

func (a *adapter) ListKeyResultsById(id []int64) ([]*model.KeyResult, error) {
	var keyResult []*model.KeyResult
	err := a.db.Where("id IN ?", id).Order("id DESC").Find(&keyResult).Error
	if err != nil {
		return nil, err
	}
	return keyResult, nil
}

func (a *adapter) ListKeyResultsByObjectiveId(objectiveId int64) ([]*model.KeyResult, error) {
	var keyResult []*model.KeyResult
	err := a.db.Where("objective_id = ?", objectiveId).Order("id DESC").Find(&keyResult).Error
	if err != nil {
		return nil, err
	}
	return keyResult, nil
}

func (a *adapter) CreateKeyResult(keyResult *model.KeyResult) (int64, error) {
	locker.Mux.Lock()
	defer locker.Mux.Unlock()

	// sequence
	sequence := int64(0)
	var max model.KeyResult
	err := a.db.Where("`uid` = ? AND `topic` = ?", keyResult.Uid, keyResult.Topic).Order("sequence DESC").Take(&max).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if max.Sequence > 0 {
		sequence = max.Sequence
	}
	sequence += 1

	keyResult.Sequence = sequence
	err = a.db.Create(&keyResult).Error
	if err != nil {
		return 0, err
	}

	// init value record
	if keyResult.CurrentValue > 0 {
		err = a.db.Create(&model.KeyResultValue{
			KeyResultId: keyResult.Id,
			Value:       keyResult.CurrentValue,
		}).Error
		if err != nil {
			return 0, err
		}
	}

	return keyResult.Id, nil
}

func (a *adapter) UpdateKeyResult(keyResult *model.KeyResult) error {
	return a.db.Model(&model.KeyResult{}).
		Where("`uid` = ? AND `topic` = ? AND sequence = ?", keyResult.Uid, keyResult.Topic, keyResult.Sequence).
		UpdateColumns(map[string]interface{}{
			"title":        keyResult.Title,
			"memo":         keyResult.Memo,
			"target_value": keyResult.TargetValue,
			"value_mode":   keyResult.ValueMode,
		}).Error
}

func (a *adapter) DeleteKeyResult(id int64) error {
	return a.db.Where("id = ?", id).Delete(&model.KeyResult{}).Error
}

func (a *adapter) DeleteKeyResultBySequence(uid types.Uid, topic string, sequence int64) error {
	return a.db.Where("`uid` = ? AND `topic` = ? AND sequence = ?", uid.UserId(), topic, sequence).Delete(&model.KeyResult{}).Error
}

func (a *adapter) AggregateObjectiveValue(id int64) error {
	result := model.KeyResult{}
	err := a.db.Model(&model.KeyResult{}).Where("objective_id = ?", id).
		Select("SUM(current_value) as current_value, SUM(target_value) as target_value").Take(&result).Error
	if err != nil {
		return err
	}
	return a.db.Model(&model.Objective{}).Where("id = ?", id).UpdateColumns(map[string]interface{}{
		"current_value": result.CurrentValue,
		"total_value":   result.TargetValue,
	}).Error
}

func (a *adapter) AggregateKeyResultValue(id int64) error {
	keyResult, err := a.GetKeyResultByID(id)
	if err != nil {
		return err
	}
	var value sql.NullString // fixme
	switch keyResult.ValueMode {
	case model.ValueSumMode:
		err = a.db.Model(&model.KeyResultValue{}).Where("key_result_id = ?", id).
			Select("SUM(`value`) as `value`").Pluck("value", &value).Error
	case model.ValueLastMode:
		err = a.db.Model(&model.KeyResultValue{}).Where("key_result_id = ?", id).
			Order("created_at DESC").Limit(1).Pluck("value", &value).Error
	case model.ValueAvgMode:
		err = a.db.Model(&model.KeyResultValue{}).Where("key_result_id = ?", id).
			Select("AVG(`value`) as `value`").Pluck("value", &value).Error
	case model.ValueMaxMode:
		err = a.db.Model(&model.KeyResultValue{}).Where("key_result_id = ?", id).
			Select("MAX(`value`) as `value`").Pluck("value", &value).Error
	}
	if err != nil {
		return err
	}

	currentValue, _ := strconv.ParseInt(value.String, 10, 64)
	return a.db.Model(&model.KeyResult{}).Where("id = ?", id).UpdateColumns(map[string]interface{}{
		"current_value": currentValue,
	}).Error
}

func (a *adapter) CreateKeyResultValue(keyResultValue *model.KeyResultValue) (int64, error) {
	err := a.db.Create(&keyResultValue).Error
	if err != nil {
		return 0, err
	}
	return keyResultValue.Id, nil
}

func (a *adapter) GetKeyResultValues(keyResultId int64) ([]*model.KeyResultValue, error) {
	var values []*model.KeyResultValue
	err := a.db.Where("key_result_id = ?", keyResultId).Order("id DESC").Find(&values).Error
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (a *adapter) CreateTodo(todo *model.Todo) (int64, error) {
	locker.Mux.Lock()
	defer locker.Mux.Unlock()

	// sequence
	sequence := int64(0)
	var max model.Todo
	err := a.db.Where("`uid` = ? AND `topic` = ?", todo.Uid, todo.Topic).Order("sequence DESC").Take(&max).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if max.Sequence > 0 {
		sequence = max.Sequence
	}
	sequence += 1

	todo.Sequence = sequence
	err = a.db.Create(&todo).Error
	if err != nil {
		return 0, nil
	}
	return todo.Id, nil
}

func (a *adapter) ListTodos(uid types.Uid, topic string) ([]*model.Todo, error) {
	var items []*model.Todo
	err := a.db.
		Where("`uid` = ? AND `topic` = ?", uid.UserId(), topic).
		Order("priority DESC").
		Order("created_at DESC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (a *adapter) ListRemindTodos(uid types.Uid, topic string) ([]*model.Todo, error) {
	var items []*model.Todo
	err := a.db.
		Where("`uid` = ? AND `topic` = ?", uid.UserId(), topic).
		Where("complete <> ?", 1).
		Where("is_remind_at_time = ?", 1).
		Order("priority DESC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (a *adapter) GetTodo(id int64) (*model.Todo, error) {
	var find model.Todo
	err := a.db.Where("id = ?", id).First(&find).Error
	if err != nil {
		return nil, err
	}
	return &find, nil
}

func (a *adapter) GetTodoBySequence(uid types.Uid, topic string, sequence int64) (*model.Todo, error) {
	var find model.Todo
	err := a.db.
		Where("`uid` = ? AND `topic` = ? AND sequence = ?", uid.UserId(), topic, sequence).
		First(&find).Error
	if err != nil {
		return nil, err
	}
	return &find, nil
}

func (a *adapter) CompleteTodo(id int64) error {
	return a.db.Model(&model.Todo{}).
		Where("id = ?", id).
		Update("complete", true).Error
}

func (a *adapter) CompleteTodoBySequence(uid types.Uid, topic string, sequence int64) error {
	return a.db.Model(&model.Todo{}).
		Where("`uid` = ? AND `topic` = ? AND sequence = ?", uid.UserId(), topic, sequence).
		Update("complete", true).Error
}

func (a *adapter) UpdateTodo(todo *model.Todo) error {
	return a.db.Model(&model.Todo{}).
		Where("`uid` = ? AND `topic` = ? AND sequence = ?", todo.Uid, todo.Topic, todo.Sequence).
		UpdateColumns(map[string]interface{}{
			"content":  todo.Content,
			"category": todo.Category,
			"remark":   todo.Remark,
			"priority": todo.Priority,
		}).Error
}

func (a *adapter) DeleteTodo(id int64) error {
	return a.db.Where("id = ?", id).Delete(&model.Todo{}).Error
}

func (a *adapter) DeleteTodoBySequence(uid types.Uid, topic string, sequence int64) error {
	return a.db.
		Where("`uid` = ? AND `topic` = ? AND sequence = ?", uid.UserId(), topic, sequence).
		Delete(&model.Todo{}).Error
}

func init() {
	store.RegisterAdapter(&adapter{})
}
