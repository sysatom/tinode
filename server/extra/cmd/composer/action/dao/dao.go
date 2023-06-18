package dao

import (
	"encoding/json"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/jsonco"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"os"
)

func GenerationAction(c *cli.Context) error {
	conffile := c.String("config")

	file, err := os.Open(conffile)
	if err != nil {
		panic(err)
	}

	config := configType{}
	jr := jsonco.New(file)
	if err = json.NewDecoder(jr).Decode(&config); err != nil {
		panic(err)
	}

	if config.StoreConfig.UseAdapter != "mysql" {
		panic("error adapter")
	}
	if config.StoreConfig.Adapters.Mysql.DSN == "" {
		panic("error adapter dsn")
	}
	dsn := config.StoreConfig.Adapters.Mysql.DSN

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./server/extra/store/dao",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery,
	})

	g.UseDB(db)

	g.ApplyBasic(
		model.User{},
		model.Topic{},
		model.Message{},
		model.Credential{},
	)
	g.ApplyBasic(
		model.Config{},
		model.OAuth{},
		model.Form{},
		model.Action{},
		model.Session{},
		model.Page{},
		model.Data{},
		model.Url{},
		model.Behavior{},
		model.Instruct{},
		model.Workflow{},
		model.Parameter{},
	)
	g.ApplyBasic(
		model.Objective{},
		model.KeyResult{},
		model.KeyResultValue{},
		model.Todo{},
		model.Counter{},
	)
	g.Execute()

	return nil
}

type configType struct {
	StoreConfig struct {
		UseAdapter string `json:"use_adapter"`
		Adapters   struct {
			Mysql struct {
				DSN string `json:"dsn"`
			} `json:"mysql"`
		} `json:"adapters"`
	} `json:"store_config"`
}
