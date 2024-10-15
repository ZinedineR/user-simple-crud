package migration

import (
	"user-simple-crud/internal/entity"
	"user-simple-crud/pkg/database"
)

func AutoMigration(CpmDB *database.Database) {
	CpmDB.MigrateDB(

		&entity.User{})
	//&entity.SMSLog{}
}
