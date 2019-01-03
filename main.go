package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// sqlite3|postgres
	dbType = "sqlite3"

	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "afkc"

	sqliteDbPath = "./test.db"
)

type Task struct {
	ID          string `json:"task_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Estimate    int    `json:"estimate,omitempty"`
	Spent       int    `json:"spent,omitempty"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", GetTasks).Methods("GET")
	router.HandleFunc("/tasks", CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", GetTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	db := connectToDb()
	defer db.Close()

	tasks := make([]Task, 0)

	rows, err := db.Query("SELECT * FROM tasks")
	panicOnError(err)
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID,
			&task.Title,
			&task.Description,
			&task.Estimate,
			&task.Spent)
		panicOnError(err)

		tasks = append(tasks, task)
	}
	defer rows.Close()

	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)

	newID, err := uuid.NewRandom()
	panicOnError(err)
	task.ID = newID.String()

	db := connectToDb()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO tasks VALUES (?, ?, ?, ?, ?)")
	panicOnError(err)
	_, err = stmt.Exec(task.ID, task.Title, task.Description, task.Estimate, task.Spent)
	panicOnError(err)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	params := mux.Vars(r)

	db := connectToDb()
	defer db.Close()

	err := db.QueryRow("SELECT * FROM tasks WHERE id=?", params["id"]).
		Scan(&task.ID, &task.Title, &task.Description, &task.Estimate, &task.Spent)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	panicOnError(err)

	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	db := connectToDb()
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM tasks where id=?")
	panicOnError(err)

	res, err := stmt.Exec(params["id"])
	panicOnError(err)

	affected, err := res.RowsAffected()
	panicOnError(err)
	if affected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if affected != 1 {
		panic("Unexpected DB response during deletion, deleted more than 1 row")
	}

	json.NewEncoder(w).Encode(params["id"])
}

func getConnectionString(db string) (dbInfo string, success bool) {
	if db == "postgres" {
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		), true
	} else if db == "sqlite3" {
		return sqliteDbPath, true
	}
	return "", false
}

func connectToDb() *sql.DB {
	dbInfo, success := getConnectionString(dbType)
	if !success {
		panic(fmt.Sprintf("Unknown dbType: %s", dbType))
	}
	db, err := sql.Open(dbType, dbInfo)
	panicOnError(err)
	return db
}
