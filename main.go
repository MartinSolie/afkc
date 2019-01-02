package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = ""
	dbname = "afkc"
)

type Task struct {
	ID			string	`json:"task_id,omitempty"`
	Title		string	`json:"title,omitempty"`
	Description	string	`json:"description,omitempty"`
	Estimate	int		`json:"estimate,omitempty"`
	Spent		int		`json:"spent,omitempty"`
}

var tasks []Task

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", GetTasks).Methods("GET")
	router.HandleFunc("/tasks", CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", GetTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	db := connectToDb()
	defer db.Close()

	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range tasks {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			break
		}
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	var deleted Task
	params := mux.Vars(r)
	for index, item := range tasks {
		if item.ID == params["id"] {
			deleted = item
			tasks = append(tasks[:index], tasks[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(deleted)
}

func getConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
						host, port, user, password, dbname,
	)
}

func connectToDb() *sql.DB {
	psqlInfo := getConnectionString()
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}
