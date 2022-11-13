//go:build mysql
// +build mysql

package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/tinode/chat/server/db/mysql"
	"github.com/tinode/chat/server/extra/locker"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/store/types"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
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

func (a *adapter) GetObjectiveByID(ctx context.Context, id int64) (*model.Objective, error) {
	var objective model.Objective
	err := a.db.WithContext(ctx).Where("id = ?", id).First(&objective).Error
	if err != nil {
		return nil, err
	}
	return &objective, nil
}

func (a *adapter) GetObjectiveBySequence(ctx context.Context, userId, sequence int64) (*model.Objective, error) {
	var objective model.Objective
	err := a.db.WithContext(ctx).Where("user_id = ? AND sequence = ?", userId, sequence).First(&objective).Error
	if err != nil {
		return nil, err
	}
	return &objective, nil
}

func (a *adapter) ListObjectives(ctx context.Context, userId int64) ([]*model.Objective, error) {
	var objectives []*model.Objective
	err := a.db.WithContext(ctx).Where("user_id = ?", userId).Order("id DESC").Find(&objectives).Error
	if err != nil {
		return nil, err
	}
	return objectives, nil
}

func (a *adapter) CreateObjective(ctx context.Context, objective *model.Objective) (int64, error) {
	locker.Mux.Lock()
	defer locker.Mux.Unlock()

	// sequence
	sequence := int64(0)
	var max model.Objective
	err := a.db.WithContext(ctx).Where("user_id = ?", objective.UserId).Order("sequence DESC").Take(&max).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if max.Sequence > 0 {
		sequence = max.Sequence
	}
	sequence += 1

	objective.Sequence = sequence
	objective.CreatedAt = time.Now().Unix()
	objective.UpdatedAt = time.Now().Unix()
	err = a.db.WithContext(ctx).Create(&objective).Error
	if err != nil {
		return 0, err
	}
	return objective.Id, nil
}

func (a *adapter) UpdateObjective(ctx context.Context, objective *model.Objective) error {
	return a.db.WithContext(ctx).Model(&model.Objective{}).
		Where("user_id = ? AND sequence = ?", objective.UserId, objective.Sequence).
		UpdateColumns(map[string]interface{}{
			"title":       objective.Title,
			"memo":        objective.Memo,
			"motive":      objective.Motive,
			"feasibility": objective.Feasibility,
			"is_plan":     objective.IsPlan,
			"plan_start":  objective.PlanStart,
			"plan_end":    objective.PlanEnd,
			"updated_at":  time.Now().Unix(),
		}).Error
}

func (a *adapter) DeleteObjective(ctx context.Context, id int64) error {
	return a.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Objective{}).Error
}

func (a *adapter) DeleteObjectiveBySequence(ctx context.Context, userId, sequence int64) error {
	return a.db.WithContext(ctx).Where("user_id = ? AND sequence = ?", userId, sequence).Delete(&model.Objective{}).Error
}

func (a *adapter) GetKeyResultByID(ctx context.Context, id int64) (*model.KeyResult, error) {
	var keyResult model.KeyResult
	err := a.db.WithContext(ctx).Where("id = ?", id).First(&keyResult).Error
	if err != nil {
		return nil, err
	}
	return &keyResult, nil
}

func (a *adapter) GetKeyResultBySequence(ctx context.Context, userId, sequence int64) (*model.KeyResult, error) {
	var keyResult model.KeyResult
	err := a.db.WithContext(ctx).Where("user_id = ? AND sequence = ?", userId, sequence).First(&keyResult).Error
	if err != nil {
		return nil, err
	}
	return &keyResult, nil
}

func (a *adapter) ListKeyResults(ctx context.Context, userId int64) ([]*model.KeyResult, error) {
	var keyResult []*model.KeyResult
	err := a.db.WithContext(ctx).Where("user_id = ?", userId).Order("id DESC").Find(&keyResult).Error
	if err != nil {
		return nil, err
	}
	return keyResult, nil
}

func (a *adapter) ListKeyResultsById(ctx context.Context, id []int64) ([]*model.KeyResult, error) {
	var keyResult []*model.KeyResult
	err := a.db.WithContext(ctx).Where("id IN ?", id).Order("id DESC").Find(&keyResult).Error
	if err != nil {
		return nil, err
	}
	return keyResult, nil
}

func (a *adapter) ListKeyResultsByObjectiveId(ctx context.Context, objectiveId int64) ([]*model.KeyResult, error) {
	var keyResult []*model.KeyResult
	err := a.db.WithContext(ctx).Where("objective_id = ?", objectiveId).Order("id DESC").Find(&keyResult).Error
	if err != nil {
		return nil, err
	}
	return keyResult, nil
}

func (a *adapter) CreateKeyResult(ctx context.Context, keyResult *model.KeyResult) (int64, error) {
	locker.Mux.Lock()
	defer locker.Mux.Unlock()

	// sequence
	sequence := int64(0)
	var max model.KeyResult
	err := a.db.WithContext(ctx).Where("user_id = ?", keyResult.UserId).Order("sequence DESC").Take(&max).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if max.Sequence > 0 {
		sequence = max.Sequence
	}
	sequence += 1

	keyResult.Sequence = sequence
	keyResult.CreatedAt = time.Now().Unix()
	keyResult.UpdatedAt = time.Now().Unix()
	err = a.db.WithContext(ctx).Create(&keyResult).Error
	if err != nil {
		return 0, err
	}

	// init value record
	if keyResult.CurrentValue > 0 {
		err = a.db.WithContext(ctx).Create(&model.KeyResultValue{
			KeyResultId: keyResult.Id,
			Value:       keyResult.CurrentValue,
			CreatedAt:   time.Now().Unix(),
		}).Error
		if err != nil {
			return 0, err
		}
	}

	return keyResult.Id, nil
}

func (a *adapter) UpdateKeyResult(ctx context.Context, keyResult *model.KeyResult) error {
	return a.db.WithContext(ctx).Model(&model.KeyResult{}).
		Where("user_id = ? AND sequence = ?", keyResult.UserId, keyResult.Sequence).
		UpdateColumns(map[string]interface{}{
			"title":        keyResult.Title,
			"memo":         keyResult.Memo,
			"target_value": keyResult.TargetValue,
			"value_mode":   keyResult.ValueMode,
			"updated_at":   time.Now().Unix(),
		}).Error
}

func (a *adapter) DeleteKeyResult(ctx context.Context, id int64) error {
	return a.db.WithContext(ctx).Where("id = ?", id).Delete(&model.KeyResult{}).Error
}

func (a *adapter) DeleteKeyResultBySequence(ctx context.Context, userId, sequence int64) error {
	return a.db.WithContext(ctx).Where("user_id = ? AND sequence = ?", userId, sequence).Delete(&model.KeyResult{}).Error
}

func (a *adapter) AggregateObjectiveValue(ctx context.Context, id int64) error {
	result := model.KeyResult{}
	err := a.db.WithContext(ctx).Model(&model.KeyResult{}).Where("objective_id = ?", id).
		Select("SUM(current_value) as current_value, SUM(target_value) as target_value").Take(&result).Error
	if err != nil {
		return err
	}
	return a.db.WithContext(ctx).Model(&model.Objective{}).Where("id = ?", id).UpdateColumns(map[string]interface{}{
		"current_value": result.CurrentValue,
		"total_value":   result.TargetValue,
		"updated_at":    time.Now().Unix(),
	}).Error
}

func (a *adapter) AggregateKeyResultValue(ctx context.Context, id int64) error {
	keyResult, err := a.GetKeyResultByID(ctx, id)
	if err != nil {
		return err
	}
	var value sql.NullInt64
	switch keyResult.ValueMode {
	case model.ValueSumMode:
		err = a.db.WithContext(ctx).Model(&model.KeyResultValue{}).Where("key_result_id = ?", id).
			Select("SUM(value) as value").Pluck("value", &value).Error
	case model.ValueLastMode:
		err = a.db.WithContext(ctx).Model(&model.KeyResultValue{}).Where("key_result_id = ?", id).
			Order("created_at DESC").Limit(1).Pluck("value", &value).Error
	case model.ValueAvgMode:
		err = a.db.WithContext(ctx).Model(&model.KeyResultValue{}).Where("key_result_id = ?", id).
			Select("AVG(value) as value").Pluck("value", &value).Error
	case model.ValueMaxMode:
		err = a.db.WithContext(ctx).Model(&model.KeyResultValue{}).Where("key_result_id = ?", id).
			Select("MAX(value) as value").Pluck("value", &value).Error
	}
	if err != nil {
		return err
	}

	return a.db.WithContext(ctx).Model(&model.KeyResult{}).Where("id = ?", id).UpdateColumns(map[string]interface{}{
		"current_value": value.Int64,
		"updated_at":    time.Now().Unix(),
	}).Error
}

func (a *adapter) CreateKeyResultValue(ctx context.Context, keyResultValue *model.KeyResultValue) (int64, error) {
	keyResultValue.CreatedAt = time.Now().Unix()
	err := a.db.WithContext(ctx).Create(&keyResultValue).Error
	if err != nil {
		return 0, err
	}
	return keyResultValue.Id, nil
}

func (a *adapter) GetKeyResultValues(ctx context.Context, keyResultId int64) ([]*model.KeyResultValue, error) {
	var values []*model.KeyResultValue
	err := a.db.WithContext(ctx).Where("key_result_id = ?", keyResultId).Order("id DESC").Find(&values).Error
	if err != nil {
		return nil, err
	}
	return values, nil
}

func init() {
	store.RegisterAdapter(&adapter{})
}
