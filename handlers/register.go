package handlers

import (
	"github.com/sirupsen/logrus"
	"html/template"
	"log"
	"main.go/DatabaseInit"
	"main.go/HashPasswordFunc"
	"main.go/structures"
	"net/http"
	"strings"
)

func Register(w http.ResponseWriter, r *http.Request) {
	DatabaseInit.DatabaseInit()
	n, err := template.ParseFiles("htmlt/registration.html")
	if err != nil {
		log.Fatal(err)
	}
	data := map[string]interface{}{}
	if r.Method == "POST" {
		f1, f2 := true, true
		data["Error"] = "You create your account succsesfull"
		var User structures.User
		User.Email = r.FormValue("useremail")
		parts := strings.Split(User.Email, "@")

		domain := parts[1]
		if domain != "gmail.com" {
			data["Error"] = "you dont use gmail.com"
			f1 = false
		}
		User.Name, User.Password = r.FormValue("username"), r.FormValue("userpassword")
		if len(User.Password) < 10 {
			data["ErrorPassword"] = "You use less then 10 symbols in your password"
		} else {
			HashedPassword, err := HashPasswordFunc.HashPassword(User.Password)
			if err != nil {
				logrus.Print(err)
			}
			_, err = DatabaseInit.Db.Exec("INSERT INTO login (name, password, email) VALUES (?, ?, ?)", User.Name, HashedPassword, User.Email)
			if err != nil {
				if strings.Contains(err.Error(), "Duplicate entry") {
					data["Error"] = "This person was registered later"
					f2 = false
				}
			}
			if f1 && f2 {
				log.Println("User was Registred")
				var id string
				err = DatabaseInit.Db.QueryRow("SELECT id FROM login WHERE name = ?", User.Name).Scan(&id)
				if err != nil {
					log.Print(err)
				}
				http.SetCookie(w, &http.Cookie{
					Name:     "user_id",
					Value:    id,
					Path:     "/",
					HttpOnly: true,
					MaxAge:   7200,
				})
				http.Redirect(w, r, "/main", http.StatusSeeOther)
			}
		}
	}
	err = n.ExecuteTemplate(w, "register", data)
	if err != nil {
		log.Fatal(err)
	}
}
