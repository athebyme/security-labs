package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var db *sql.DB
var secretKey = "super_secret_key_12345"
var dbPassword = "root:password123@tcp(localhost:3306)/mydb"

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	query := fmt.Sprintf("SELECT id FROM users WHERE username='%s' AND password='%s'", username, password)
	db.QueryRow(query)
	fmt.Fprintf(w, "Welcome, %s!", username)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	cmd := exec.Command("sh", "-c", "ping -c 1 "+host)
	output, _ := cmd.CombinedOutput()
	w.Write(output)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	fmt.Fprintf(w, "<html><body><h1>Results for: %s</h1></body></html>", q)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	data, _ := os.ReadFile("./uploads/" + name)
	w.Write(data)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	buf := make([]byte, 4096)
	n, _ := resp.Body.Read(buf)
	w.Write(buf[:n])
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	tmpl := template.Must(template.New("p").Parse("<h1>{{.}}</h1>"))
	tmpl.Execute(w, template.HTML(name))
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/file", fileHandler)
	http.HandleFunc("/proxy", proxyHandler)
	http.HandleFunc("/profile", profileHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
