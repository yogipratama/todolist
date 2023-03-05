package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Todo struct {
	Id          int
	IndexNumber int
	Name        string
	Status      string
}

func GetConnection() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/todo")
	if err != nil {
		panic(err)
	}
	return db
}

var tmpl = template.Must(template.ParseFiles("./templates/index.gohtml"))

// func get all data todo
func IndexAndInsert(writer http.ResponseWriter, request *http.Request) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	if request.Method == "POST" {
		name := request.FormValue("name")
		querySQL := "INSERT INTO todo(name) VALUES(?)"
		ctx := context.Background()
		statement, err := db.PrepareContext(ctx, querySQL)
		if err != nil {
			panic(err)
		}
		defer statement.Close()
		statement.ExecContext(ctx, name)
		http.Redirect(writer, request, "/", http.StatusMovedPermanently)
	}

	querySQL := "SELECT id, name, status FROM todo ORDER BY status DESC"
	rows, err := db.QueryContext(ctx, querySQL)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// get data todo using loop
	todo := Todo{}
	response := []Todo{}

	for rows.Next() {
		var id int
		var name, status string
		err := rows.Scan(&id, &name, &status)
		if err != nil {
			panic(err)
		}
		todo.Id = id
		todo.Name = name
		todo.Status = status
		todo.IndexNumber++
		response = append(response, todo)
	}
	tmpl.ExecuteTemplate(writer, "index.gohtml", response)
}

// func insert data todo
func Delete(writer http.ResponseWriter, request *http.Request) {
	db := GetConnection()
	defer db.Close()

	getId := request.URL.Query().Get("id")
	ctx := context.Background()
	querySQL := "DELETE FROM todo WHERE id = ?"
	statement, err := db.PrepareContext(ctx, querySQL)
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	statement.ExecContext(ctx, getId)
	http.Redirect(writer, request, "/", http.StatusMovedPermanently)
}

// func edit status todo to done
func EditStatus(writer http.ResponseWriter, request *http.Request) {
	db := GetConnection()
	defer db.Close()

	getId := request.URL.Query().Get("id")
	var status string = "1"
	ctx := context.Background()
	querySQL := "UPDATE todo SET status = ? WHERE id = ?"
	statement, err := db.PrepareContext(ctx, querySQL)
	if err != nil {
		panic(err)
	}
	defer statement.Close()
	statement.ExecContext(ctx, status, getId)
	http.Redirect(writer, request, "/", http.StatusMovedPermanently)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", IndexAndInsert)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", EditStatus)
	http.ListenAndServe(":8080", nil)
}
