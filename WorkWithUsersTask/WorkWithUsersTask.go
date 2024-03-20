package WorkWithUsersTask

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/user")
	if err != nil {
		log.Fatal(err)
	} else {
		logrus.Print("Conected with DataBase")
	}
}

func WorkWithTask(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("action")
	switch action {
	case "delete":
		vars := mux.Vars(r)
		IdOfTask := vars["id"]

		smrt, err := db.Prepare("DELETE FROM utasks WHERE id=?")
		if err != nil {
			logrus.Print(err)
		}
		_, err = smrt.Exec(IdOfTask)
		if err != nil {
			logrus.Print(err)
		}
		http.Redirect(w, r, "/main", http.StatusSeeOther)
	case "update":
		NewName := r.FormValue("NameOfnewTask")
		NewDescription := r.FormValue("descriptionOfNewTask")
		vars := mux.Vars(r)
		IdOfTask := vars["id"]
		smrt, err := db.Prepare("UPDATE utasks SET NameOfTask = ?, DescriptionOfTask = ? WHERE id = ?")
		if err != nil {
			logrus.Print(err)
		}
		_, err = smrt.Exec(NewName, NewDescription, IdOfTask)
		if err != nil {
			logrus.Print(err)
		}
		http.Redirect(w, r, "/main", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/main", http.StatusSeeOther)
	}
}
