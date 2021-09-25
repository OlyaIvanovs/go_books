package main

import (
	"database/sql"
	"encoding/json"
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

type SearchResult struct {
	Title  string
	Author string
	Year   string
	ID     string
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

	http.HandleFunc("/search", func(rw http.ResponseWriter, r *http.Request) {
		results := []SearchResult{
			SearchResult{"Moby-Dick", "Herman Melwille", "1851", "222"},
			SearchResult{"The adventures of Finn", "Mark Twain", "1884", "4444"},
			SearchResult{"The catcher in the Rye", "JD Salinger", "1951", "3333"},
		}

		encoder := json.NewEncoder(rw)
		if err := encoder.Encode(results); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Print(http.ListenAndServe(":8090", nil))
}
