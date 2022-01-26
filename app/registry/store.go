package registry

import (
	"github.com/gotoeveryone/auth-api/app/config"
	"github.com/gotoeveryone/auth-api/app/infrastructure/database"
)

// InitDatastore is initialize datastore
func InitDatastore(debug bool, db config.DB) error {
	return database.Init(debug, db)
}
