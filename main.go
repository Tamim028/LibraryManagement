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


//get all books with authors name
func allBooks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(books) //even if there are no book then return empty array
	w.WriteHeader(http.StatusOK) // responding with 200 ok
}

//get a book by ID
func singleBookByID(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") //for formatting purpose
	vars := mux.Vars(r)

	var error bool
	error = true

	//loop through books array and search for the requested book
	for _, book := range books {
		if book.ID == vars["id"] {
			error = false
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	//if there is an error the show http status
	if error == true {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Resource not found"))
		return
	}

	json.NewEncoder(w).Encode(&Book{})
	w.WriteHeader(http.StatusOK) //no error then send http ok status
}

func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book) //decode into json format

	idCounter++

	book.ID = strconv.Itoa(idCounter) //covert unique idCounter to string
	books = append(books, book)

	json.NewEncoder(w).Encode(book)
	w.WriteHeader(http.StatusOK)
}

func loanRequest(w http.ResponseWriter, r *http.Request)  {
	//do loan request and allow book loan
}

func authorization(w http.ResponseWriter, r *http.Request)  {
	//restrict user actions according to their access limit
}

//update book using book ID
func updateBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var error bool
	error = true

	var tempBook Book

	//tempBook.ID = strconv.Itoa(tempBook.ID)

	for index, item := range books {
		if item.ID == vars["id"] { //book id finding

			//book ID found so, no error generated
			error = false

			_ = json.NewDecoder(r.Body).Decode(&tempBook)

			//updating with current book data
			books[index] = tempBook
			books[index].ID = item.ID

			json.NewEncoder(w).Encode(books[index])
			return
		}
	}

	if error == true{
		//error generated, show error status and custom response
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
			tempBook = books[index]
			books = append(books[:index], books[index+1:]...) //append left slice of index and right slice of index into one
			_ = json.NewEncoder(w).Encode(tempBook)
			return
		}
	}

	if error == true{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Delete Request"))
		return
	}

	//decrease a idCounter after deleting one
	idCounter--
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Book{})
}

func homePage(w http.ResponseWriter, r *http.Request)  {
	//welcoming to homepage
	fmt.Fprint(w, "Welcome to homepage")
	w.WriteHeader(http.StatusOK)
	//log.Println(w, "Homepage endpoint hit")
}


func main() {
	//declaring mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	// Manually setting data
	books = append(books, Book{ID: "1", ISBN: "438227", Title: "Book First", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", ISBN: "454555", Title: "Book Second", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})
	books = append(books, Book{ID: "3", ISBN: "4545568", Title: "Book Third", Author: &Author{Firstname: "Malik", Lastname: "Khan"}})

	//initial idCounter setup
	idCounter = 3

	//hadling all routes
	myRouter.HandleFunc("/library", homePage)
	myRouter.HandleFunc("/library/books", allBooks).Methods("GET")
	myRouter.HandleFunc("/library/book/{id}", singleBookByID).Methods("GET")
	myRouter.HandleFunc("/library/books", addBook).Methods("POST")
	myRouter.HandleFunc("/library/books/{id}", updateBook).Methods("PUT")
	myRouter.HandleFunc("/library/books/{id}", deleteBook).Methods("DELETE")

	//hosting at 8081, used for http request
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}
