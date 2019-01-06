package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var err error
var db *sql.DB

func DB() *sql.DB {

	conn := "user:password@tcp(dburl:3306)/dbname"

	db, err = sql.Open("mysql", conn)
	if err != nil {
		fmt.Println(err)
	}

	return db

}
