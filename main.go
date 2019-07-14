package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Article type
type Article struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

// Error type is for defining error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Articles is array that holds Article type
var Articles []Article

func return404(w http.ResponseWriter, r *http.Request, m string) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(
		Error{Code: 404, Message: m},
	)
}

func set204NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}

func set201Created(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

func setContentJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	setContentJSON(w)

	// get PUT Content
	r.ParseForm()
	title := r.FormValue("title")
	desc := r.FormValue("desc")
	content := r.FormValue("content")

	key := mux.Vars(r)["id"]

	for i := range Articles {
		// Check if ID is same
		if Articles[i].ID == key {
			// Check if Form values has value if nil omit
			if title != "" {
				Articles[i].Title = title
			}

			if desc != "" {
				Articles[i].Desc = desc
			}

			if content != "" {
				Articles[i].Content = content
			}

			json.NewEncoder(w).Encode(Articles[i])
			return
		}
	}

	return404(w, r, "Article not found")
}

func singleArticle(w http.ResponseWriter, r *http.Request) {
	setContentJSON(w)

	vars := mux.Vars(r)
	key := vars["id"]

	for _, article := range Articles {
		if article.ID == key {
			json.NewEncoder(w).Encode(article)
			return
		}
	}

	return404(w, r, "Article not found")
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for i := range Articles {
		if Articles[i].ID == key {
			set204NoContent(w)
			Articles = Articles[:i+copy(Articles[i:], Articles[i+1:])]

			return
		}
	}

	return404(w, r, "Article not found")
}

func allArticles(w http.ResponseWriter, r *http.Request) {
	setContentJSON(w)

	json.NewEncoder(w).Encode(Articles)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	setContentJSON(w)
	set201Created(w)

	// To get the form data into key, value pair
	// we use r.ParseForm()
	r.ParseForm()
	var article Article

	newID := len(Articles) + 1
	article.ID = strconv.Itoa(newID)

	article.Title = r.FormValue("title")
	article.Desc = r.FormValue("desc")
	article.Content = r.FormValue("content")

	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(article)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage")
	fmt.Println("Endpoint hit: homePage")
}

func handleRequest() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/articles", allArticles).Methods("GET")

	// Routing order is important
	router.HandleFunc("/articles", createNewArticle).Methods("POST")
	router.HandleFunc("/articles/{id}", singleArticle).Methods("GET")
	router.HandleFunc("/articles/{id}", updateArticle).Methods("PATCH")
	router.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")

	http.Handle("/", router) // Necessary to setup default routing to mux
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func main() {
	Articles = []Article{
		Article{ID: "1", Title: "Test 1", Desc: "Desc", Content: "Content"},
		Article{ID: "2", Title: "Test 2", Desc: "Desc 2", Content: "Content 2"},
		Article{ID: "3", Title: "Test 3", Desc: "Desc 3", Content: "Content 3"},
	}

	handleRequest()
}
