package rest

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ilhamsyahids/bookshelf-template/storage"
	"github.com/ilhamsyahids/bookshelf-template/utils"
)

type API struct {
	bookStorage storage.Storage
}

type APIConfig struct {
	BookStorage storage.Storage
}

func NewAPI(config APIConfig) (*API, error) {
	return &API{bookStorage: config.BookStorage}, nil
}

func (api *API) GetHandler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", api.serveHealthCheck)

	r.Get("/books", api.serveGetBooks)

	r.Get("/books/:id", api.serveGetBookByID)

	r.Post("/books", api.serveCreateBook)

	r.Put("/books/:id", api.serveUpdateBook)

	r.Delete("/books/:id", api.serveDeleteBook)

	return r
}

func (api *API) serveHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It's working!"))
}

func (api *API) serveGetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := api.bookStorage.GetBooks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// output success response
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(utils.NewSuccessResp(books))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(buf.Bytes())
}

func (api *API) serveGetBookByID(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")

	book, err := api.bookStorage.GetBookByID(bookID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	// Output success response with the book
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(utils.NewSuccessResp(book))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (api *API) serveCreateBook(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming request body
	var newBook storage.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Generate a new ID for the book
	newBook.ID = uuid.New().String()

	// Add the new book to the storage
	err = api.bookStorage.AddBook(newBook)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Return success response with the created book
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(utils.NewSuccessResp(newBook))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

func (api *API) serveUpdateBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")

	// Parse the request body
	updateRequest := &storage.Book{}
	if err := json.NewDecoder(r.Body).Decode(updateRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Retrieve the book from storage
	book, err := api.bookStorage.GetBookByID(bookID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	// Update the book fields if the corresponding fields in the update request are not empty
	if updateRequest.ISBN != "" {
		book.ISBN = updateRequest.ISBN
	}
	if updateRequest.Title != "" {
		book.Title = updateRequest.Title
	}
	if updateRequest.Author != "" {
		book.Author = updateRequest.Author
	}
	if updateRequest.Published != "" {
		book.Published = updateRequest.Published
	}

	// Update the book in storage
	// (assuming you have a method `UpdateBookByID` in the `storage` package)
	err = api.bookStorage.UpdateBookByID(bookID, book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Output success response
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(utils.NewSuccessResp(book))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (api *API) serveDeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")

	err := api.bookStorage.DeleteBookByID(bookID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
