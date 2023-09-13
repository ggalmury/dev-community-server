package initializers

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
)

var DB *gorm.DB

func DBConnection() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true},
	})

	if err != nil {
		panic("Failed to connect to database")
	}
}
