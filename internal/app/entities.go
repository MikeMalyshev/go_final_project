package app

type Storage interface {
	AddTask(task Task) (int64, error)
	GetTaskByID(id string) (Task, error)
	GetTaskList(searchString string, maxLen int64) ([]Task, error)
	UpdateTask(task Task) error
	RemoveTask(id string) error
	FindTask(title, date string) (string, error)
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskList struct {
	List []Task `json:"tasks"`
}
