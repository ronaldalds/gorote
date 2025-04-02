package settings

import (
	"github.com/ronaldalds/base-go-api/internal/config/databases"
	"github.com/ronaldalds/base-go-api/internal/config/envs"
)

func Config() error {
	envs.Load()
	databases.LoadSQL()
	databases.LoadRedis()
	// databases.LoadNOSQL()
	if err := Ready(); err != nil {
		return err
	}
	return nil
}
