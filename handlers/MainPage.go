package handlers

import (
	"github.com/sirupsen/logrus"
	"html/template"
	"main.go/DatabaseInit"
	"main.go/structures"
	"net/http"
	"sort"
	"time"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	data := prepareData()
	n, err := template.ParseFiles("htmlt/main.html")
	if err != nil {
		logrus.Print(err)
	}

	switch r.Method {
	case "POST":
		handlePostRequest(w, r, data)
	case "GET":
		handleGetRequest(w, r, data)
	}
	err = n.ExecuteTemplate(w, "main", data)
	if err != nil {
		logrus.Print(err)
	}
}

func prepareData() map[string]interface{} {
	data := map[string]interface{}{}
	data["Premission"] = true
	return data
}

func handlePostRequest(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		data["Premission"] = false
		data["Error"] = "Login to have premission to add task"
	} else {
		handleTaskCreation(w, r, cookie)
	}
}

func handleTaskCreation(w http.ResponseWriter, r *http.Request, cookie *http.Cookie) {
	var PostTask structures.UserTasks
	PostTask.TimeOfTask, PostTask.NameOfTask, PostTask.DescriptioOfTask = r.FormValue("UserTimeOfTask"), r.FormValue("NameOfUserTask"), r.FormValue("UserDescriptionOfTask")
	if PostTask.TimeOfTask == "" {
		insertTaskWithoutTime(PostTask, cookie)
	} else {
		insertTaskWithTime(PostTask, cookie)
	}
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

func insertTaskWithoutTime(PostTask structures.UserTasks, cookie *http.Cookie) {
	smrt, err := DatabaseInit.Db.Prepare("INSERT INTO utasks (NameOfTask, DescriptionOfTask, User_id) VALUES (?, ?, ?)")
	if err != nil {
		logrus.Print(err)
	}
	defer smrt.Close()
	userID := cookie.Value
	_, err = smrt.Exec(PostTask.NameOfTask, PostTask.DescriptioOfTask, userID)
	if err != nil {
		logrus.Print(err)
	}
}

func insertTaskWithTime(PostTask structures.UserTasks, cookie *http.Cookie) {
	smrt, err := DatabaseInit.Db.Prepare("INSERT INTO utasks (NameOfTask, DescriptionOfTask, User_id, DateOfTask ) VALUES (?, ?, ?, ?)")
	if err != nil {
		logrus.Print(err)
	}
	defer smrt.Close()
	userID := cookie.Value
	_, err = smrt.Exec(PostTask.NameOfTask, PostTask.DescriptioOfTask, userID, PostTask.TimeOfTask)
	if err != nil {
		logrus.Print(err)
	}
}

func handleGetRequest(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	cookie, _ := r.Cookie("user_id")
	if cookie == nil {
		data["Premission"] = false
		data["Error"] = "Please login to have a premission to your tasks "
	} else {
		data["Tasks"] = fetchTasks(cookie)
	}
}

func fetchTasks(cookie *http.Cookie) []structures.Tasks {
	var Task []structures.Tasks
	timeSlice := fetchTaskTimes(cookie)
	for _, i := range timeSlice {
		Task = append(Task, fetchTasksWithTime(i, cookie)...)
	}
	Task = append(Task, fetchTasksWithoutTime(cookie)...)
	return Task
}

func fetchTaskTimes(cookie *http.Cookie) []time.Time {
	DataRow, _ := DatabaseInit.Db.Query("SELECT DateOfTask FROM utasks WHERE User_id = ? AND DateOfTask IS NOT NULL ", cookie.Value)
	IsOk := DataRow.Next()
	var (
		Time      string
		timeSlice []time.Time
	)
	for ok := IsOk; ok; ok = DataRow.Next() {
		err := DataRow.Scan(&Time)
		if err != nil {
			logrus.Print(err)
		}
		data, err := time.Parse("2006-01-02 15:04:05", Time)
		if err != nil {
			logrus.Print(err)
		}
		timeSlice = append(timeSlice, data)
	}
	sort.Slice(timeSlice, func(i, j int) bool {
		return timeSlice[i].Before(timeSlice[j])
	})
	return timeSlice
}

func fetchTasksWithTime(i time.Time, cookie *http.Cookie) []structures.Tasks {
	var Task []structures.Tasks
	var oneTask structures.UserTasks
	row, _ := DatabaseInit.Db.Query("SELECT id ,NameOfTask, DescriptionOfTask, DateOfTask  FROM  utasks WHERE User_id = ? AND DateOfTask = ?", cookie.Value, i)
	isOk := row.Next()
	if isOk {
		for ok := isOk; ok; ok = row.Next() {
			err := row.Scan(&oneTask.ID, &oneTask.NameOfTask, &oneTask.DescriptioOfTask, &oneTask.TimeOfTask)
			if err != nil {
				Task = append(Task, structures.Tasks{Id: oneTask.ID, Name: oneTask.NameOfTask, Description: oneTask.DescriptioOfTask})
			} else {
				Task = append(Task, structures.Tasks{Id: oneTask.ID, Name: oneTask.NameOfTask, Description: oneTask.DescriptioOfTask, Data: oneTask.TimeOfTask})
			}
		}
	}

	return Task
}

func fetchTasksWithoutTime(cookie *http.Cookie) []structures.Tasks {
	var Task []structures.Tasks
	var oneTask structures.UserTasks
	row, _ := DatabaseInit.Db.Query("SELECT id ,NameOfTask, DescriptionOfTask  FROM  utasks WHERE User_id = ? AND DateOfTask IS NULL ", cookie.Value)
	IsOk := row.Next()
	if IsOk {
		for ok := IsOk; ok; ok = row.Next() {
			err := row.Scan(&oneTask.ID, &oneTask.NameOfTask, &oneTask.DescriptioOfTask)
			if err != nil {
				logrus.Print(err)
			}
			Task = append(Task, structures.Tasks{Id: oneTask.ID, Name: oneTask.NameOfTask, Description: oneTask.DescriptioOfTask})
		}
	}
	return Task
}
