package sqlhandler

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SqlHandler struct {
	DB *gorm.DB
}

func NewSqlHandler() *SqlHandler {
	sqlHandler := new(SqlHandler)
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	dbname := os.Getenv("MYSQL_DATABASE")
	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, host, dbname)
	DB, err := gorm.Open(mysql.Open(connection), &gorm.Config{})
	if err != nil {
		log.Fatalln(connection + "database can't connect")
	}
	sqlHandler.DB = DB
	return sqlHandler
}
