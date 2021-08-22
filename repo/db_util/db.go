package db_util

import (
	"fmt"

	"eos-arb/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDb(name, pass string) *gorm.DB {

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(localhost:3306)/eos_arb?charset=utf8&parseTime=True&loc=Local", name, pass)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(
		&model.PairIdList{},
	)
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(10)
	return db
}
