package handlers

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"main.go/DatabaseInit"
	"main.go/structures"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	DatabaseInit.DatabaseInit()
	n, err := template.ParseFiles("htmlt/login.html")
	if err != nil {
		log.Print(err)
	}
	data := map[string]interface{}{}

	if r.Method == "POST" {
		var User structures.User
		smrt, err := DatabaseInit.Db.Prepare("SELECT password FROM login WhERE name = ?")
		if err != nil {
			log.Print(err)
		}
		defer smrt.Close()
		User.Name = r.FormValue("username")
		User.Password = r.FormValue("userpassword")
		var password string
		err = smrt.QueryRow(User.Name).Scan(&password)

		if err != nil {
			if err == sql.ErrNoRows {
				data["Error"] = "User with that login is undefiend"
			} else {
				logrus.Print(err)
			}
		} else {
			err = bcrypt.CompareHashAndPassword([]byte(password), []byte(User.Password))
			if err == nil {
				err = DatabaseInit.Db.QueryRow("SELECT id FROM login WHERE name = ?", User.Name).Scan(&User.ID)
				if err != nil {
					log.Fatal(err)
				}
				http.SetCookie(w, &http.Cookie{
					Name:     "user_id",
					Value:    User.ID,
					Path:     "/",
					HttpOnly: true,
					MaxAge:   7200,
				})
				http.Redirect(w, r, "/main", http.StatusSeeOther)
			} else {
				data["Error"] = "Your information is ancorrect"
			}
		}
	}
	err = n.ExecuteTemplate(w, "login", data)
	if err != nil {
		log.Print(err)
	}
}
