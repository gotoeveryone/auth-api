package registry

import (
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/infrastructure/database"
)

// InitDatastore is initialize datastore
func InitDatastore() error {
	return database.Init(config.AppConfig.App.Debug, config.AppConfig.DB)
}
