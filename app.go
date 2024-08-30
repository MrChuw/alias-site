package main

import (
    "database/sql"
    _ "fmt"
    "net/http"
    "html/template"
    "log"
    _ "github.com/mattn/go-sqlite3"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
)

var db *sql.DB

func initDB() {
    var err error
    db, err = sql.Open("sqlite3", "./db/db.sqlite3")
    if err != nil {
        log.Fatal(err)
    }
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS markdown_tables (
        id TEXT PRIMARY KEY,
        content TEXT,
        table_name TEXT
    )`)
    if err != nil {
        log.Fatal(err)
    }
}

func submitTable(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        markdownContent := r.FormValue("markdown")
        tableName := r.FormValue("table_name")
        if markdownContent != "" && tableName != "" {
			tableID := uuid.New().String()
            _, err := db.Exec("INSERT INTO markdown_tables (id, content, table_name) VALUES (?, ?, ?)", tableID, markdownContent, tableName)
            if err != nil {
                http.Error(w, "Failed to insert data", http.StatusInternalServerError)
                return
            }
			log.Println("tableID Created: ", tableID)
            http.Redirect(w, r, "/table/"+tableID, http.StatusSeeOther)
            return
        }
        http.Error(w, "No Markdown content or table name provided", http.StatusBadRequest)
        return
    }
    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

func viewTable(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tableID := vars["table_id"]
	log.Println("tableID: ", tableID)
    row := db.QueryRow("SELECT content, table_name FROM markdown_tables WHERE id = ?", tableID)
    var content, tableName string
    if err := row.Scan(&content, &tableName); err != nil {
        http.Error(w, "Table not found", http.StatusNotFound)
        return
    }
    tmpl, err := template.ParseFiles("templates/table.html")
    if err != nil {
        http.Error(w, "Template parsing error", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, map[string]string{
        "Content":   content,
        "TableName": tableName,
    })
}

func index(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "templates/index.html")
}

func main() {
    initDB()
    r := mux.NewRouter()

    // Serve static files
    fs := http.FileServer(http.Dir("templates"))
    r.Handle("/favicon.ico", http.FileServer(http.Dir("templates")))
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

    r.HandleFunc("/submit", submitTable).Methods(http.MethodPost)
    r.HandleFunc("/table/{table_id:[a-zA-Z0-9-]+}", viewTable)
    r.HandleFunc("/", index)

	log.Println("Server starting on port 80...")
    log.Fatal(http.ListenAndServe(":80", r))
}
