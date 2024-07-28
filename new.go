package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Todo struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	IsCompleted bool   `json:"isCompleted"`
}

var todos = []Todo{
	{Id: 1, Name: "Learn Go", IsCompleted: false},
	{Id: 2, Name: "Learn Alpine", IsCompleted: false},
	{Id: 3, Name: "Go to the gym", IsCompleted: true},
}

var templates map[string]*template.Template

// Load templates on program initialization
func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["index.html"] = template.Must(template.ParseFiles("index.html"))
	templates["todo.html"] = template.Must(template.ParseFiles("todo.html"))
}

// handlers
func submitTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	isCompleted := r.FormValue("isCompleted") == "true"
	todo := Todo{Id: len(todos) + 1, Name: name, IsCompleted: isCompleted}
	todos = append(todos, todo)

	err := templates["todo.html"].ExecuteTemplate(w, "todo.html", todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(todos)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to marshal todos", http.StatusInternalServerError)
		return
	}

	tmplData := map[string]interface{}{"Todos": string(data)}
	err = templates["index.html"].ExecuteTemplate(w, "index.html", tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit-todo/", submitTodoHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
