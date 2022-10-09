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
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `key` = ?", int64(uid), topic, key).First(&find).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if find.ID > 0 {
		return a.db.
			Model(&model.Config{}).
			Where("`uid` = ? AND `topic` = ? AND `key` = ?", int64(uid), topic, key).
			Update("value", value).Error
	} else {
		return a.db.Create(&model.Config{
			Uid:   int64(uid),
			Topic: topic,
			Key:   key,
			Value: value,
		}).Error
	}
}

func (a *adapter) ConfigGet(uid types.Uid, topic, key string) (model.JSON, error) {
	var find model.Config
	err := a.db.Where("`uid` = ? AND `topic` = ? AND `key` = ?", int(uid), topic, key).First(&find).Error
	if err != nil {
		return nil, err
	}
	return find.Value, nil
}

func init() {
	store.RegisterAdapter(&adapter{})
}
