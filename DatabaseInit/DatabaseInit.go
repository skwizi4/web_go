package DatabaseInit

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"log"
)

var Db *sql.DB

func DatabaseInit() {
	var Err error
	Db, Err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/user")
	if Err != nil {
		log.Fatal(Err)
	} else {
		logrus.Print("Conected with DataBase")
	}
}
