//go:build mysql
// +build mysql

package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/tinode/chat/server/db/mysql"
	"github.com/tinode/chat/server/extra/pkg/flog"
	"github.com/tinode/chat/server/extra/pkg/locker"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/dao"
	"github.com/tinode/chat/server/extra/store/model"
	extraTypes "github.com/tinode/chat/server/extra/types"
	serverStore "github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

const adapterName = "mysql"

type adapter struct {
	db *gorm.DB
}

func (a *adapter) Open() error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Warn,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                   // Disable color
		},
	)
	db, err := gorm.Open(mysqlDriver.New(mysqlDriver.Config{Conn: mysql.RawDB}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	// Initialize dao
	dao.SetDefault(db)
	// gorm
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

func wildcard(str string) string {
	return fmt.Sprintf("%%%s%%", str)
}

func (a *adapter) GetBotUsers() ([]*model.User, error) {
	var find []*model.User
	err := a.db.
		Table("users").
		Select("id, createdat, access, lastseen, useragent, public->'$.fn' AS fn, trusted->'$.verified' AS verified").
		Where("trusted->'$.verified' = TRUE AND JSON_CONTAINS(tags, '[\"bot\"]')").
		Find(&find).Error
	if err != nil {
		return nil, err
	}
	return find, nil
}

func (a *adapter) GetNormalUsers() ([]*model.User, error) {
	var find []*model.User
	err := a.db.
		Table("users").
		Select("id, createdat, access, lastseen, useragent, public->'$.fn' AS fn, trusted->'$.verified' AS verified").
		Where("lastseen IS NOT NULL").
		Find(&find).Error
	if err != nil {
		return nil, err
	}
	return find, nil
}

func (a *adapter) GetGroupTopics(owner types.Uid) ([]*model.Topic, error) {
	var find []*model.Topic
	err := a.db.
		Table("topics").
		Select("id, createdat, state, touchedat, name, owner, access, seqid, delid, public, trusted, public->'$.fn' AS fn, trusted->'$.verified' AS verified").
		Where("`owner` = ? AND `name` LIKE 'grp%' AND `usebt` = false", serverStore.DecodeUid(owner)).
		Find(&find).Error
	if err != nil {
		return nil, err
	}
	return find, nil
}

func (a *adapter) SearchMessages(uid types.Uid, searchTopic string, filter string) ([]*model.Message, error) {
	var find []*model.Message
	err := a.db.
		Table("messages").
		Select("id, `from`, topic, content->'$.txt' AS txt, content as raw, seqid, createdat").
		Where("`from` = ? AND (`content` LIKE ? OR `content`->'$.txt' LIKE ?) AND `topic` <> ?", serverStore.DecodeUid(uid), wildcard(filter), wildcard(filter), searchTopic).
		Order("createdat DESC").
		Limit(10).
		Find(&find).Error
	if err != nil {
		return nil, err
	}
	return find, nil
}

func (a *adapter) GetMessage(topic string, seqId int) (model.Message, error) {
	var find model.Message
	err := a.db.Where("`topic` = ? AND `seqid` = ?", topic, seqId).First(&find).Error
	if err != nil {
		return model.Message{}, err
	}
	return find, nil
}

func (a *adapter) GetCredentials() ([]*model.Credential, error) {
	var find []*model.Credential
	err := a.db.Find(&find).Error
	if err != nil {
		return nil, err
	}
	return find, nil
}

func (a *adapter) DataSet(uid types.Uid, topic, key string, value extraTypes.KV) error {
	var find model.Data
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `key` = ?", uid.UserId(), topic, key).First(&find).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if find.ID > 0 {
		return a.db.
			Model(&model.Data{}).
			Where("`uid` = ? AND `topic` = ? AND `key` = ?", uid.UserId(), topic, key).
			Update("value", value).Error
	} else {
		return a.db.Create(&model.Data{
			UID:   uid.UserId(),
			Topic: topic,
			Key:   key,
			Value: model.JSON(value),
		}).Error
	}
}

func (a *adapter) DataGet(uid types.Uid, topic, key string) (extraTypes.KV, error) {
	var find model.Data
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `key` = ?", uid.UserId(), topic, key).First(&find).Error
	if err != nil {
		return nil, err
	}
	return extraTypes.KV(find.Value), nil
}

func (a *adapter) DataList(uid types.Uid, topic string, filter extraTypes.DataFilter) ([]*model.Data, error) {
	var list []*model.Data
	builder := a.db.Where("`uid` = ? AND `topic` = ?", uid.UserId(), topic)
	if filter.Prefix != nil {
		builder = builder.Where("`key` LIKE ?", fmt.Sprintf("%s%%", *filter.Prefix))
	}
	if filter.CreatedStart != nil {
		builder = builder.Where("created_at >= ?", filter.CreatedStart)
	}
	if filter.CreatedEnd != nil {
		builder = builder.Where("created_at <= ?", filter.CreatedEnd)
	}
	err := builder.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (a *adapter) DataDelete(uid types.Uid, topic string, key string) error {
	return a.db.Where("`uid` = ? AND `topic` = ? AND `key` = ?", uid.UserId(), topic, key).Delete(&model.Data{}).Error
}

func (a *adapter) ConfigSet(uid types.Uid, topic, key string, value extraTypes.KV) error {
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
			UID:   uid.UserId(),
			Topic: topic,
			Key:   key,
			Value: model.JSON(value),
		}).Error
	}
}

func (a *adapter) ConfigGet(uid types.Uid, topic, key string) (extraTypes.KV, error) {
	var find model.Config
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `key` = ?", uid.UserId(), topic, key).First(&find).Error
	if err != nil {
		return nil, err
	}
	return extraTypes.KV(find.Value), nil
}

func (a *adapter) OAuthSet(oauth model.OAuth) error {
	var find model.OAuth
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `type` = ?", oauth.UID, oauth.Topic, oauth.Type).First(&find).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if find.ID > 0 {
		return a.db.
			Model(&model.OAuth{}).
			Where("`uid` = ? AND `topic` = ? AND `type` = ?", oauth.UID, oauth.Topic, oauth.Type).
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
			FormID: formId,
			UID:    form.UID,
			Topic:  form.Topic,
			Schema: form.Schema,
			Values: form.Values,
			State:  form.State,
			Extra:  form.Extra,
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

func (a *adapter) ActionSet(topic string, seqId int, action model.Action) error {
	var find model.Action
	err := a.db.Where("`topic` = ? AND `seqid` = ?", topic, seqId).First(&find).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if find.ID > 0 {
		return a.db.
			Model(&model.Action{}).
			Where("`topic` = ? AND `seqid` = ?", topic, seqId).
			Updates(map[string]interface{}{
				"value": action.Value,
				"state": action.State,
			}).Error
	} else {
		return a.db.Create(&model.Action{
			UID:   action.UID,
			Topic: topic,
			Seqid: int32(seqId),
			Value: action.Value,
			State: action.State,
		}).Error
	}
}

func (a *adapter) ActionGet(topic string, seqId int) (model.Action, error) {
	var find model.Action
	err := a.db.Where("`topic` = ? AND `seqid` = ?", topic, seqId).First(&find).Error
	if err != nil {
		return model.Action{}, err
	}
	return find, nil
}

func (a *adapter) SessionCreate(session model.Session) error {
	return a.db.Create(&model.Session{
		UID:    session.UID,
		Topic:  session.Topic,
		RuleID: session.RuleID,
		Init:   session.Init,
		Values: session.Values,
		State:  session.State,
	}).Error
}

func (a *adapter) SessionSet(uid types.Uid, topic string, session model.Session) error {
	var find model.Session
	err := a.db.Where("`uid` = ? AND `topic` = ?", uid.UserId(), topic).Order("created_at DESC").First(&find).Error
	if err != nil {
		return err
	}
	return a.db.
		Model(&model.Session{}).
		Where("`id` = ?", find.ID).
		Updates(map[string]interface{}{
			"values": session.Values,
		}).Error
}

func (a *adapter) SessionState(uid types.Uid, topic string, state model.SessionState) error {
	var find model.Session
	err := a.db.Where("`uid` = ? AND `topic` = ?", uid.UserId(), topic).Order("created_at DESC").First(&find).Error
	if err != nil {
		return err
	}
	return a.db.
		Model(&model.Session{}).
		Where("`id` = ?", find.ID).
		Updates(map[string]interface{}{
			"state": state,
		}).Error
}

func (a *adapter) SessionGet(uid types.Uid, topic string) (model.Session, error) {
	var find model.Session
	err := a.db.Where("`uid` = ? AND `topic` = ?", uid.UserId(), topic).Order("created_at DESC").First(&find).Error
	if err != nil {
		return model.Session{}, err
	}
	return find, nil
}

func (a *adapter) PipelineCreate(pipeline model.Pipeline) error {
	return a.db.Create(&model.Pipeline{
		UID:     pipeline.UID,
		Topic:   pipeline.Topic,
		Flag:    pipeline.Flag,
		RuleID:  pipeline.RuleID,
		Version: pipeline.Version,
		Stage:   pipeline.Stage,
		Values:  pipeline.Values,
		State:   pipeline.State,
	}).Error
}

func (a *adapter) PipelineState(uid types.Uid, topic string, pipeline model.Pipeline) error {
	return a.db.
		Model(&model.Pipeline{}).
		Where("`uid` = ? AND `topic` = ? AND `flag` = ?", uid.UserId(), topic, pipeline.Flag).
		Updates(map[string]interface{}{
			"state": pipeline.State,
		}).Error
}

func (a *adapter) PipelineStep(uid types.Uid, topic string, pipeline model.Pipeline) error {
	return a.db.
		Model(&model.Pipeline{}).
		Where("`uid` = ? AND `topic` = ? AND `flag` = ?", uid.UserId(), topic, pipeline.Flag).
		Updates(map[string]interface{}{
			"stage": pipeline.Stage,
		}).Error
}

func (a *adapter) PipelineGet(uid types.Uid, topic string, flag string) (model.Pipeline, error) {
	var find model.Pipeline
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `flag` = ?", uid.UserId(), topic, flag).Order("created_at DESC").First(&find).Error
	if err != nil {
		return model.Pipeline{}, err
	}
	return find, nil
}

func (a *adapter) BehaviorSet(behavior model.Behavior) error {
	return a.db.Create(&behavior).Error
}

func (a *adapter) BehaviorGet(uid types.Uid, flag string) (model.Behavior, error) {
	var find model.Behavior
	err := a.db.Where("`uid` = ? AND `flag` = ?", uid.UserId(), flag).First(&find).Error
	if err != nil {
		return model.Behavior{}, err
	}
	return find, nil
}

func (a *adapter) BehaviorList(uid types.Uid) ([]*model.Behavior, error) {
	var list []*model.Behavior
	err := a.db.Where("`uid` = ?", uid.UserId()).Order("id DESC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (a *adapter) BehaviorIncrease(uid types.Uid, flag string, number int) error {
	return a.db.
		Model(&model.Behavior{}).
		Where("`uid` = ? AND `flag` = ?", uid.UserId(), flag).
		UpdateColumn("count", gorm.Expr("count + ?", number)).Error
}

func (a *adapter) ParameterSet(flag string, params extraTypes.KV, expiredAt time.Time) error {
	var find model.Parameter
	err := a.db.Where("`flag` = ?", flag).First(&find).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if find.ID > 0 {
		return a.db.
			Model(&model.Parameter{}).
			Where("`flag` = ?", flag).
			Update("expired_at", expiredAt).Error
	} else {
		return a.db.Create(&model.Parameter{
			Flag:      flag,
			Params:    model.JSON(params),
			ExpiredAt: expiredAt,
		}).Error
	}
}

func (a *adapter) ParameterGet(flag string) (model.Parameter, error) {
	var find model.Parameter
	err := a.db.Where("`flag` = ?", flag).First(&find).Error
	if err != nil {
		return model.Parameter{}, err
	}
	return find, nil
}

func (a *adapter) ParameterDelete(flag string) error {
	_, err := dao.Q.Parameter.Where(dao.Parameter.Flag.Eq(flag)).Delete()
	return err
}

func (a *adapter) UrlCreate(url model.Url) error {
	return a.db.Create(&model.Url{
		Flag:  url.Flag,
		URL:   url.URL,
		State: url.State,
	}).Error
}

func (a *adapter) UrlViewIncrease(flag string) error {
	return a.db.
		Model(&model.Url{}).
		Where("`flag` = ?", flag).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (a *adapter) UrlState(flag string, state model.UrlState) error {
	return a.db.
		Model(&model.Url{}).
		Where("`flag` = ?", flag).
		Updates(map[string]interface{}{
			"state": state,
		}).Error
}

func (a *adapter) UrlGetByFlag(flag string) (model.Url, error) {
	var find model.Url
	err := a.db.Where("`flag` = ?", flag).First(&find).Error
	if err != nil {
		return model.Url{}, err
	}
	return find, nil
}

func (a *adapter) UrlGetByUrl(url string) (model.Url, error) {
	var find model.Url
	err := a.db.Where("`url` = ?", url).First(&find).Error
	if err != nil {
		return model.Url{}, err
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

func (a *adapter) CreateInstruct(instruct *model.Instruct) (int64, error) {
	if instruct.ExpireAt.Before(time.Now()) {
		return 0, errors.New("expire time error")
	}
	err := a.db.Create(&instruct)
	if err != nil {
		return 0, nil
	}
	return int64(instruct.ID), nil
}

func (a *adapter) ListInstruct(uid types.Uid, isExpire bool) ([]*model.Instruct, error) {
	var items []*model.Instruct
	builder := a.db.
		Where("`uid` = ?", uid.UserId()).
		Where("state = ?", model.InstructCreate)
	if isExpire {
		builder.Where("expire_at < ?", time.Now())
	} else {
		builder.Where("expire_at >= ?", time.Now())
	}

	err := builder.
		Order("priority DESC").
		Order("updated_at DESC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (a *adapter) UpdateInstruct(instruct *model.Instruct) error {
	return a.db.Model(&model.Todo{}).
		Where("`no` = ?", instruct.No).
		UpdateColumns(map[string]interface{}{
			"state": instruct.State,
		}).Error
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
	ctx := context.Background()
	l := locker.NewLocker()
	lock, err := l.Acquire(ctx, fmt.Sprintf("chatbot:objective:create:%s", objective.UID), 10*time.Second)
	if err != nil {
		return 0, err
	}
	defer func(lock *redislock.Lock, ctx context.Context) {
		_ = lock.Release(ctx)
	}(lock, ctx)

	// sequence
	sequence := int32(0)
	var max model.Objective
	err = a.db.Where("`uid` = ? AND `topic` = ?", objective.UID, objective.Topic).Order("sequence DESC").Take(&max).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if max.Sequence > 0 {
		sequence = max.Sequence
	}
	sequence += 1

	if objective.CreatedData.IsZero() {
		objective.CreatedData = time.Now()
	}
	if objective.UpdatedDate.IsZero() {
		objective.UpdatedDate = time.Now()
	}
	objective.Sequence = sequence
	err = a.db.Create(&objective).Error
	if err != nil {
		return 0, err
	}
	return int64(objective.ID), nil
}

func (a *adapter) UpdateObjective(objective *model.Objective) error {
	return a.db.Model(&model.Objective{}).
		Where("`uid` = ? AND `topic` = ? AND sequence = ?", objective.UID, objective.Topic, objective.Sequence).
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
	ctx := context.Background()
	l := locker.NewLocker()
	lock, err := l.Acquire(ctx, fmt.Sprintf("chatbot:key_result:create:%s", keyResult.UID), 10*time.Second)
	if err != nil {
		return 0, err
	}
	defer func(lock *redislock.Lock, ctx context.Context) {
		_ = lock.Release(ctx)
	}(lock, ctx)

	// sequence
	sequence := int32(0)
	var max model.KeyResult
	err = a.db.Where("`uid` = ? AND `topic` = ?", keyResult.UID, keyResult.Topic).Order("sequence DESC").Take(&max).Error
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
			KeyResultID: keyResult.ID,
			Value:       keyResult.CurrentValue,
		}).Error
		if err != nil {
			return 0, err
		}
	}

	return int64(keyResult.ID), nil
}

func (a *adapter) UpdateKeyResult(keyResult *model.KeyResult) error {
	return a.db.Model(&model.KeyResult{}).
		Where("`uid` = ? AND `topic` = ? AND sequence = ?", keyResult.UID, keyResult.Topic, keyResult.Sequence).
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
	var value sql.NullInt64
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

	return a.db.Model(&model.KeyResult{}).Where("id = ?", id).UpdateColumns(map[string]interface{}{
		"current_value": value.Int64,
	}).Error
}

func (a *adapter) CreateKeyResultValue(keyResultValue *model.KeyResultValue) (int64, error) {
	err := a.db.Create(&keyResultValue).Error
	if err != nil {
		return 0, err
	}
	return int64(keyResultValue.ID), nil
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
	ctx := context.Background()
	l := locker.NewLocker()
	lock, err := l.Acquire(ctx, fmt.Sprintf("chatbot:todo:create:%s", todo.UID), 10*time.Second)
	if err != nil {
		return 0, err
	}
	defer func(lock *redislock.Lock, ctx context.Context) {
		_ = lock.Release(ctx)
	}(lock, ctx)

	// sequence
	sequence := int32(0)
	var max model.Todo
	err = a.db.Where("`uid` = ? AND `topic` = ?", todo.UID, todo.Topic).Order("sequence DESC").Take(&max).Error
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
	return int64(todo.ID), nil
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
		Where("`uid` = ? AND `topic` = ? AND sequence = ?", todo.UID, todo.Topic, todo.Sequence).
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

func (a *adapter) CreateReview(review *model.Review) (int64, error) {
	q := dao.Q.Review
	err := q.Create(review)
	if err != nil {
		return 0, err
	}
	return int64(review.ID), nil
}

func (a *adapter) UpdateReview(review *model.Review) error {
	q := dao.Q.Review
	_, err := q.Updates(review)
	return err
}

func (a *adapter) ListReviews(uid types.Uid, topic string) ([]*model.Review, error) {
	q := dao.Q.Review
	return q.Where(q.UID.Eq(uid.UserId()), q.Topic.Eq(topic)).
		Order(q.UpdatedAt.Desc()).Find()
}

func (a *adapter) GetReviewByID(id int64) (*model.Review, error) {
	q := dao.Q.Review
	return q.Where(q.ID.Eq(int32(id))).First()
}

func (a *adapter) CreateReviewEvaluation(evaluation *model.ReviewEvaluation) (int64, error) {
	q := dao.Q.ReviewEvaluation
	err := q.Create(evaluation)
	if err != nil {
		return 0, err
	}
	return int64(evaluation.ID), nil
}

func (a *adapter) UpdateReviewEvaluation(evaluation *model.ReviewEvaluation) error {
	q := dao.Q.ReviewEvaluation
	_, err := q.Updates(evaluation)
	return err
}

func (a *adapter) ListReviewEvaluations(uid types.Uid, topic string, reviewID int64) ([]*model.ReviewEvaluation, error) {
	q := dao.Q.ReviewEvaluation
	return q.Where(q.UID.Eq(uid.UserId()), q.Topic.Eq(topic), q.ReviewID.Eq(int32(reviewID))).
		Order(q.UpdatedAt.Desc()).Find()
}

func (a *adapter) GetReviewEvaluationByID(id int64) (*model.ReviewEvaluation, error) {
	q := dao.Q.ReviewEvaluation
	return q.Where(q.ID.Eq(int32(id))).First()
}

func (a *adapter) CreateCycle(cycle *model.Cycle) (int64, error) {
	q := dao.Q.Cycle
	err := q.Create(cycle)
	if err != nil {
		return 0, err
	}
	return int64(cycle.ID), nil
}

func (a *adapter) UpdateCycle(cycle *model.Cycle) error {
	q := dao.Q.Cycle
	_, err := q.Updates(cycle)
	return err
}

func (a *adapter) ListCycles(uid types.Uid, topic string) ([]*model.Cycle, error) {
	q := dao.Q.Cycle
	return q.Where(q.UID.Eq(uid.UserId()), q.Topic.Eq(topic)).
		Order(q.UpdatedAt.Desc()).Find()
}

func (a *adapter) GetCycleByID(id int64) (*model.Cycle, error) {
	q := dao.Q.Cycle
	return q.Where(q.ID.Eq(int32(id))).First()
}

func (a *adapter) CreateCounter(counter *model.Counter) (int64, error) {
	err := a.db.Create(&counter)
	if err != nil {
		return 0, nil
	}
	a.record(int64(counter.ID), counter.Digit)
	return int64(counter.ID), nil
}

func (a *adapter) IncreaseCounter(id, amount int64) error {
	err := a.db.Model(&model.Counter{}).
		Where("id = ?", id).
		Update("digit", gorm.Expr("digit + ?", amount)).Error
	if err != nil {
		return err
	}
	a.record(id, amount)
	return nil
}

func (a *adapter) DecreaseCounter(id, amount int64) error {
	err := a.db.Model(&model.Counter{}).
		Where("id = ?", id).
		Update("digit", gorm.Expr("digit - ?", amount)).Error
	if err != nil {
		return err
	}
	a.record(id, -amount)
	return nil
}

func (a *adapter) ListCounter(uid types.Uid, topic string) ([]*model.Counter, error) {
	var items []*model.Counter
	err := a.db.Where("`uid` = ? AND `topic` = ?", uid.UserId(), topic).
		Order("updated_at DESC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (a *adapter) record(id, digit int64) {
	err := a.db.Exec("INSERT INTO `chatbot_counter_records` ( `counter_id`, `digit`, `created_at`) VALUES (?, ?, ?)",
		id, digit, time.Now()).Error
	if err != nil {
		flog.Error(err)
	}
}

func (a *adapter) GetCounter(id int64) (model.Counter, error) {
	var find model.Counter
	err := a.db.Where("id = ?", id).First(&find).Error
	if err != nil {
		return model.Counter{}, err
	}
	return find, nil
}

func (a *adapter) GetCounterByFlag(uid types.Uid, topic string, flag string) (model.Counter, error) {
	var find model.Counter
	err := a.db.Where("`uid` = ? AND `topic` = ? AND flag = ?", uid.UserId(), topic, flag).First(&find).Error
	if err != nil {
		return model.Counter{}, err
	}
	return find, nil
}

func (a *adapter) CreateWorkflow(workflow *model.Workflow, dag *model.Dag, triggers []*model.WorkflowTrigger) (int64, error) {
	q := dao.Q
	err := q.Transaction(func(tx *dao.Query) error {
		err := tx.Workflow.Create(workflow)
		if err != nil {
			return err
		}

		err = tx.Dag.Create(dag)
		if err != nil {
			return err
		}

		for _, trigger := range triggers {
			err = tx.WorkflowTrigger.Create(trigger)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return int64(workflow.ID), nil
}

func (a *adapter) GetWorkflow(id int64) (*model.Workflow, error) {
	q := dao.Q.Workflow
	return q.Where(q.ID.Eq(int32(id))).First()
}

func (a *adapter) UpdateWorkflowState(id int64, state model.WorkflowState) error {
	q := dao.Q.Workflow
	_, err := q.Where(q.ID.Eq(int32(id))).UpdateSimple(q.State.Value(state))
	return err
}

func (a *adapter) ListWorkflows(uid types.Uid, topic string) ([]*model.Workflow, error) {
	q := dao.Q.Workflow
	return q.Where(q.UID.Eq(uid.UserId())).Where(q.Topic.Eq(topic)).Find()
}

func (a *adapter) IncreaseWorkflowCount(id int64, successful int, failed int, running int, canceled int) error {
	q := dao.Q.Workflow
	_, err := q.Where(q.ID.Eq(int32(id))).UpdateSimple(
		q.SuccessfulCount.Value(int32(successful)),
		q.FailedCount.Value(int32(failed)),
		q.RunningCount.Value(int32(running)),
		q.CanceledCount.Value(int32(canceled)))
	return err
}

func (a *adapter) DeleteWorkflow(id int64) error {
	q := dao.Q.Workflow
	_, err := q.Where(q.ID.Eq(int32(id))).Delete()
	return err
}

func (a *adapter) GetDag(id int64) (*model.Dag, error) {
	q := dao.Q.Dag
	return q.Where(q.ID.Eq(int32(id))).First()
}

func (a *adapter) GetJob(id int64) (*model.Job, error) {
	q := dao.Q.Job
	return q.Where(q.ID.Eq(int32(id))).First()
}

func (a *adapter) DeleteJob(id int64) error {
	q := dao.Q.Job
	_, err := q.Where(q.ID.Eq(int32(id))).Delete()
	return err
}

func (a *adapter) ListJobs(workflowID int64) ([]*model.Job, error) {
	q := dao.Q.Job
	return q.Where(q.WorkflowID.Eq(int32(workflowID))).Find()
}

func (a *adapter) GetJobsByState(state model.JobState) ([]*model.Job, error) {
	q := dao.Q.Job
	return q.Where(q.State.Eq(state)).Find()
}

func (a *adapter) UpdateJobState(id int64, state model.JobState) error {
	q := dao.Q.Job
	_, err := q.Where(q.ID.Eq(int32(id))).Update(q.State, state)
	return err
}

func (a *adapter) UpdateStepState(id int64, state model.StepState) error {
	q := dao.Q.Step
	_, err := q.Where(q.ID.Eq(int32(id))).Update(q.State, state)
	return err
}

func (a *adapter) CreateStep(step *model.Step) (int64, error) {
	q := dao.Q.Step
	err := q.Create(step)
	if err != nil {
		return 0, err
	}
	return int64(step.ID), nil
}

func (a *adapter) CreateSteps(steps []*model.Step) error {
	q := dao.Q
	return q.Transaction(func(tx *dao.Query) error {
		for _, step := range steps {
			err := tx.Step.Create(step)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *adapter) GetStepsByState(state model.StepState) ([]*model.Step, error) {
	q := dao.Q.Step
	return q.Where(q.State.Eq(state)).Find()
}

func (a *adapter) GetStepsByDepend(jobId int64, depend []string) ([]*model.Step, error) {
	q := dao.Q.Step
	return q.Where(q.JobID.Eq(int32(jobId)), q.Columns(q.Depend).In(field.Values(depend))).Find()
}

func Init() {
	store.RegisterAdapter(&adapter{})
}
