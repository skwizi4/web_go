package structures

type Tasks struct {
	Name        string
	Id          int
	Data        string
	Description string
}

type UserTasks struct {
	NameOfTask       string
	DescriptioOfTask string
	ID               int
	TimeOfTask       string
	User_id          string
}

type User struct {
	ID       string
	Name     string
	Password string
	Email    string
}
