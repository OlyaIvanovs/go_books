package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	Name     string
	DBStatus bool
}

func main() {
	templates := template.Must(template.ParseFiles("templates/index.html"))

	os.Remove("dev.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating dev.db...")
	file, err := os.Create("dev.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("dev.db created")

	db, err := sql.Open("sqlite3", "/home/olga/projects/go_books/dev.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := Page{Name: "Gopher"}
		if name := r.FormValue("name"); name != "" {
			p.Name = name
		}

		err = db.Ping()
		if err != nil {
			fmt.Println("no connections", err)
		}

		p.DBStatus = db.Ping() == nil

		if err := templates.ExecuteTemplate(w, "index.html", p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Print(http.ListenAndServe(":8080", nil))
}
