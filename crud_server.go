package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
)

type Employee struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Simple CRUD server.")
}

func dbConn() *sql.DB {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "goDB"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Insert(w http.ResponseWriter, r *http.Request) {
	conn := dbConn()
	defer conn.Close()

	//curl -X POST -d "name=Andrew&age=22" http://localhost:8080/insert
	if r.Method == "POST" {
		//get name and age of employee
		name := r.FormValue("name")
		age := r.FormValue("age")
		_, err := conn.Query("INSERT INTO employees(name, age) VALUES(?,?)", name, age)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("The record has been inserted.")
	}
}

func Read(w http.ResponseWriter, r *http.Request) {
	conn := dbConn()
	defer conn.Close()

	//curl -X GET http://localhost:8080/read?id=1 or send this request with a web browser
	if r.Method == "GET" {
		//get id from http://localhost:8080/read?id=1
		id := r.URL.Query().Get("id")
		rows, err := conn.Query("SELECT id, name, age FROM employees WHERE id=?", id)
		if err != nil {
			log.Fatal(err)
		}

		employee := Employee{}
		for rows.Next() {
			err := rows.Scan(&employee.Id, &employee.Name, &employee.Age)
			if err != nil {
				log.Fatal(err)
			}
		}
		data, err := json.Marshal(employee)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(data)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	conn := dbConn()
	defer conn.Close()

	//curl -X PUT -d "id=1&name=Andrew&age=20" http://localhost:8080/update
	if r.Method == "PUT" {
		//get id, name and age of employee
		id := r.FormValue("id")
		name := r.FormValue("name")
		age := r.FormValue("age")
		_, err := conn.Query("UPDATE employees SET name=?, age=? WHERE id=?", name, age, id)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("The record has been updated.")
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	conn := dbConn()
	defer conn.Close()

	//curl -X DELETE http://localhost:8080/delete?id=1
	if r.Method == "DELETE" {
		//get id from http://localhost:8080/delete?id=1
		id := r.URL.Query().Get("id")
		_, err := conn.Query("DELETE FROM employees WHERE id=?", id)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("The record has been deleted.")
	}
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/read", Read)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
