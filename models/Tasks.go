package models

var counter = 1

type Task struct {
	ID        *int   `json:"id,string"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type UpdateReq struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func GenerateUniqueID() int {
	id := counter
	counter++
	return id
}
