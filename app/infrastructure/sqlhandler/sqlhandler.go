package sqlhandler

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SqlHandler struct {
	DB *gorm.DB
}

func NewSqlHandler() *SqlHandler {
	sqlHandler := new(SqlHandler)
	dsn := "root:deeptrack@tcp(mysql-container:3306)/DeepTrack?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(dsn + "database can't connect")
	}
	sqlHandler.DB = DB
	return sqlHandler
}
