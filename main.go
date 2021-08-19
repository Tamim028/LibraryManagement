package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var idCounter = 0

// Book struct (Model)
type Book struct {
	ID     string  `json:"id"`
	ISBN   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

//User struct (Model)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Author struct (Model)
type Author struct {
	//UserID string `json:"user_id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

//books slice
var books []Book


//get all books and their authors
func allBooks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(books) //even if there are no book then return empty array
	w.WriteHeader(http.StatusOK) // responding with 200 ok
}

//get a book by ID
func singleBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var err bool
	err = true

	for _, book := range books {
		if book.ID == vars["id"] {
			err = false
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	if err == true {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Resource not found"))
		return
	}

	json.NewEncoder(w).Encode(&Book{})
	w.WriteHeader(http.StatusOK)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	idCounter++

	book.ID = strconv.Itoa(idCounter)
	books = append(books, book)

	json.NewEncoder(w).Encode(book)
	w.WriteHeader(http.StatusOK)
}

func loadRequest(w http.ResponseWriter, r *http.Request)  {

}

func authorization(w http.ResponseWriter, r *http.Request)  {

}

func updateBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var error bool
	error = true

	var tempBook Book

	//tempBook.ID = strconv.Itoa(tempBook.ID)

	for index, item := range books {
		if item.ID == vars["id"] {

			error = false

			_ = json.NewDecoder(r.Body).Decode(&tempBook)
			books[index] = tempBook
			books[index].ID = item.ID

			json.NewEncoder(w).Encode(books[index])
			return
		}
	}

	if error == true{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Update Request"))
		return
	}

	json.NewEncoder(w).Encode(&Book{})
	w.WriteHeader(http.StatusOK)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var error bool
	error = true

	var tempBook Book
	for index, item := range books {
		if item.ID == vars["id"] {

			error = false
			//tempBook.ID = books[index].ID
			tempBook = books[index]
			books = append(books[:index], books[index+1:]...)
			_ = json.NewEncoder(w).Encode(tempBook)
			return
		}
	}

	if error == true{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Delete Request"))
		return
	}

	idCounter--
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Book{})
}

func homePage(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprint(w, "Welcome to homepage")
	log.Println(w, "Homepage endpoint hit")
}


func main() {
	myRouter := mux.NewRouter().StrictSlash(true)

	// Manually set data
	books = append(books, Book{ID: "1", ISBN: "438227", Title: "Book First", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", ISBN: "454555", Title: "Book Second", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})
	books = append(books, Book{ID: "3", ISBN: "4545568", Title: "Book Third", Author: &Author{Firstname: "Malik", Lastname: "Khan"}})

	idCounter = 3

	myRouter.HandleFunc("/library", homePage)
	myRouter.HandleFunc("/library/books", allBooks).Methods("GET")
	myRouter.HandleFunc("/library/book/{id}", singleBookByID).Methods("GET")
	myRouter.HandleFunc("/library/books", addBook).Methods("POST")
	myRouter.HandleFunc("/library/books/{id}", updateBook).Methods("PUT")
	myRouter.HandleFunc("/library/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}