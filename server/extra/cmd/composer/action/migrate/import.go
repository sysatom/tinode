package migrate

import (
	"encoding/json"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	extraMigrate "github.com/tinode/chat/server/extra/store/migrate"
	"github.com/tinode/jsonco"
	"github.com/urfave/cli/v2"
	"os"
)

func ImportAction(c *cli.Context) error {
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

	d, err := iofs.New(extraMigrate.Fs, "migrations")
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, fmt.Sprintf("mysql://%s", dsn))
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil {
		panic(err)
	}
	fmt.Println("migrate done")
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