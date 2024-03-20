package handlers

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"main.go/DatabaseInit"
	"main.go/HashPasswordFunc"
	"main.go/structures"
	"net/http"
	"strings"
)

// Handler for registration
func Register(w http.ResponseWriter, r *http.Request) {
	data := PrepareData()
	if r.Method == "POST" {
		HandlePostRequest(w, r, data)
	} else {
		HandleGetRequest(w, r, data)
	}
}
func PrepareData() map[string]interface{} {
	data := map[string]interface{}{}
	return data
}

// Handle GET requests
func HandleGetRequest(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	n, err := template.ParseFiles("htmlt/registration.html")
	if err != nil {
		log.Fatal(err)
	}
	err = n.ExecuteTemplate(w, "register", data)
	if err != nil {
		log.Fatal(err)
	}
}

// Handle POST requests
func HandlePostRequest(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	DatabaseInit.DatabaseInit()
	data["Error"] = "You create your account succsesfull"
	var User structures.User
	User.Email = r.FormValue("useremail")
	User.Name, User.Password = r.FormValue("username"), r.FormValue("userpassword")
	if !isPasswordValid(User.Password) {
		data["Error"] = "This password was already used"
		HandleGetRequest(w, r, data)
	} else {
		if isEmailValid(User.Email) {
			registerUser(User, data, w, r)
		} else {
			data["Error"] = "This email is uccorrect"
			HandleGetRequest(w, r, data)
		}
	}
}

// Validate email
func isEmailValid(email string) bool {
	parts := strings.Split(email, "@")
	domain := parts[1]
	if domain != "gmail.com" {
		return false
	}
	return true
}

// Validate password
func isPasswordValid(password string) bool {
	if len(password) < 10 {
		return false
	} else {
		row, err := DatabaseInit.Db.Query("SELECT password FROM login")
		if err != nil {
			logrus.Print(err)
		}
		IsOk := row.Next()
		var DataBasePassword []byte
		var CheckPassword bool
		for ok := IsOk; ok; ok = row.Next() {
			err = row.Scan(&DataBasePassword)
			if err != nil {
				logrus.Print(err)
			}
			err = bcrypt.CompareHashAndPassword(DataBasePassword, []byte(password))
			if err == nil {
				CheckPassword = false
			} else {
				CheckPassword = true
			}
		}
		return CheckPassword
	}
}

// Register user
func registerUser(User structures.User, data map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	HashedPassword, err := HashPasswordFunc.HashPassword(User.Password)
	if err != nil {
		logrus.Print(err)
	}
	_, err = DatabaseInit.Db.Exec("INSERT INTO login (name, password, email) VALUES (?, ?, ?)", User.Name, HashedPassword, User.Email)
	if strings.Contains(err.Error(), "Duplicate entry") {
		data["Error"] = "This person was registered later"
	} else {
		log.Println("User was Registred")
		err = DatabaseInit.Db.QueryRow("SELECT id FROM login WHERE name = ?", User.Name).Scan(&User.ID)
		if err != nil {
			log.Print(err)
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "user_id",
			Value:    User.ID,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   7200,
		})
		http.Redirect(w, r, "/main", http.StatusSeeOther)
	}
}
