package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/ilhamsyahids/bookshelf-template/storage"
	"github.com/ilhamsyahids/bookshelf-template/utils"
	"gopkg.in/validator.v2"
)

type API struct {
	bookStorage storage.Storage
}

type APIConfig struct {
	BookStorage storage.Storage `validate:"nonnil"`
}

func NewAPI(config APIConfig) (*API, error) {
	err := validator.Validate(config)
	if err != nil {
		return nil, fmt.Errorf("invalid API config: %w", err)
	}

	return &API{bookStorage: config.BookStorage}, nil
}

func (api *API) GetHandler() http.Handler {
	r := chi.NewRouter()

	r.Get("/", api.serveHealthCheck)

	r.Get("/books", api.serveGetBooks)
	r.Get("/books/{id}", api.serveGetBookByID)
	r.Post("/books", api.serveCreateBook)
	r.Put("/books/{id}", api.serveUpdateBook)
	r.Delete("/books/{id}", api.serveDeleteBook)

	return r
}

func (api *API) serveHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It's working!"))
}

// Path: GET `/books“
func (api *API) serveGetBooks(w http.ResponseWriter, r *http.Request) {
	// get query params (query, page, limit)
	query := r.URL.Query().Get("query")

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		render.Render(w, r, utils.NewErrorResp(http.StatusBadRequest, ErrInvalidPage.Error()))
		return
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		render.Render(w, r, utils.NewErrorResp(http.StatusBadRequest, ErrInvalidLimit.Error()))
		return
	}

	// get books from storage
	books, err := api.bookStorage.GetBooks(query, page, limit)
	if err != nil {
		render.Render(w, r, utils.NewErrorResp(http.StatusInternalServerError, err.Error()))
		return
	}

	// return success response
	render.Render(w, r, utils.NewSuccessResp(books))
}

type createBookReq struct {
	ISBN      string `json:"isbn" validate:"nonzero"`
	Title     string `json:"title" validate:"nonzero"`
	Author    string `json:"author" validate:"nonzero"`
	Published string `json:"published" validate:"nonzero"`
}

func (b *createBookReq) Bind(r *http.Request) error {
	if b.ISBN == "" {
		return ErrMissingISBN
	}
	if b.Title == "" {
		return ErrMissingTitle
	}
	if b.Author == "" {
		return ErrMissingAuthor
	}
	if b.Published == "" {
		return ErrMissingPublished
	}
	return nil
}

// Path: GET `/books/{id}`
func (api *API) serveGetBookByID(w http.ResponseWriter, r *http.Request) {
	// get path params (id)
	idStr := chi.URLParam(r, "id")
	// validate path params (id)
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		render.Render(w, r, utils.NewErrorResp(http.StatusBadRequest, ErrInvalidID.Error()))
		return
	}

	// get book from storage
	book, err := api.bookStorage.GetBookByID(idStr)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			render.Render(w, r, utils.NewErrorResp(http.StatusNotFound, ErrNotFoundBook.Error()))
			return
		}
		render.Render(w, r, utils.NewErrorResp(http.StatusInternalServerError, err.Error()))
		return
	}

	// return success response
	render.Render(w, r, utils.NewSuccessResp(book))
}

// Path: POST `/books`
func (api *API) serveCreateBook(w http.ResponseWriter, r *http.Request) {
	// get body request
	bodyReq := &createBookReq{}
	// validate body request
	err := render.Bind(r, bodyReq)
	if err != nil {
		render.Render(w, r, utils.NewErrorResp(http.StatusBadRequest, err.Error()))
		return
	}

	// create book
	book := storage.Book{
		ISBN:      bodyReq.ISBN,
		Title:     bodyReq.Title,
		Author:    bodyReq.Author,
		Published: bodyReq.Published,
	}
	resBook, err := api.bookStorage.CreateBook(book)
	if err != nil {
		render.Render(w, r, utils.NewErrorResp(http.StatusInternalServerError, err.Error()))
		return
	}
	book.ID = resBook.ID

	// return success response
	render.Render(w, r, utils.NewSuccessResp(book))
}

type updateBookReq struct {
	ISBN      string `json:"isbn"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Published string `json:"published"`
}

func (b *updateBookReq) Bind(r *http.Request) error {
	if b.ISBN == "" && b.Title == "" && b.Author == "" && b.Published == "" {
		return ErrMissingUpdateData
	}
	return nil
}

// Path: PUT `/books/{id}`
func (api *API) serveUpdateBook(w http.ResponseWriter, r *http.Request) {
	// get path params (id)
	idStr := chi.URLParam(r, "id")
	// validate path params (id)
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		render.Render(w, r, utils.NewErrorResp(http.StatusBadRequest, ErrInvalidID.Error()))
		return
	}

	// get body request
	bodyReq := &updateBookReq{}
	// validate body request
	err = render.Bind(r, bodyReq)
	if err != nil {
		render.Render(w, r, utils.NewErrorResp(http.StatusBadRequest, ErrMissingUpdateData.Error()))
		return
	}

	// update book from storage
	book := &storage.Book{
		ISBN:      bodyReq.ISBN,
		Title:     bodyReq.Title,
		Author:    bodyReq.Author,
		Published: bodyReq.Published,
	}
	book, err = api.bookStorage.UpdateBook(idStr, *book)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			render.Render(w, r, utils.NewErrorResp(http.StatusNotFound, ErrNotFoundBook.Error()))
			return
		}
		render.Render(w, r, utils.NewErrorResp(http.StatusInternalServerError, err.Error()))
		return
	}

	// return success response
	render.Render(w, r, utils.NewSuccessResp(book))
}

// Path: DELETE `/books/{id}`
func (api *API) serveDeleteBook(w http.ResponseWriter, r *http.Request) {
	// get path params (id)
	idStr := chi.URLParam(r, "id")
	// validate path params (id)
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		render.Render(w, r, utils.NewErrorResp(http.StatusBadRequest, ErrInvalidID.Error()))
		return
	}

	// delete book from storage
	err = api.bookStorage.DeleteBook(idStr)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			render.Render(w, r, utils.NewErrorResp(http.StatusNotFound, ErrNotFoundBook.Error()))
			return
		}
		render.Render(w, r, utils.NewErrorResp(http.StatusInternalServerError, err.Error()))
		return
	}

	// return success response
	render.Render(w, r, utils.NewSuccessResp(nil))
}
