//go:build mysql
// +build mysql

package mysql

import (
	"errors"
	"github.com/tinode/chat/server/db/mysql"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/store/types"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func (a *adapter) FormSet(uid types.Uid, topic string, seqId int, values map[string]interface{}, state int) error {
	var find model.Form
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `seqid` = ?", uid.UserId(), topic, seqId).First(&find).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if find.ID > 0 {
		return a.db.
			Model(&model.Form{}).
			Where("`uid` = ? AND `topic` = ? AND `seqid` = ?", uid.UserId(), topic, seqId).
			Updates(map[string]interface{}{"values": values, "state": state}).Error
	} else {
		return a.db.Create(&model.Form{
			Uid:    uid.UserId(),
			Topic:  topic,
			SeqId:  seqId,
			Values: values,
			State:  model.FormState(state),
		}).Error
	}
}

func (a *adapter) FormGet(uid types.Uid, topic string, seqId int) (model.Form, error) {
	var find model.Form
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `seqid` = ?", uid.UserId(), topic, seqId).First(&find).Error
	if err != nil {
		return model.Form{}, err
	}
	return find, nil
}

func init() {
	store.RegisterAdapter(&adapter{})
}
