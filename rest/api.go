package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
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

	r.Get("/books/:id", api.serveGetBooksByID)

	r.Post("/books", api.serveCreateBook)

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

type addBookReq struct {
	ISBN      string `json:"isbn" validate:"nonzero"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Published string `json:"published"`
}

func (addNew *addBookReq) Bind(r *http.Request) error {
	err := validator.Validate(addNew)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (api *API) serveCreateBook(w http.ResponseWriter, r *http.Request) {
	newReq := &addBookReq{}
	err := render.Bind(r, newReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("New Request: %v\n", newReq)
}
