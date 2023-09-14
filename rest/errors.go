package rest

import (
	"errors"
)

var (
	ErrInvalidPage  = errors.New("ERR_INVALID_PAGE")
	ErrInvalidLimit = errors.New("ERR_INVALID_LIMIT")

	ErrMissingTitle     = errors.New("ERR_TITLE_MISSING")
	ErrMissingAuthor    = errors.New("ERR_AUTHOR_MISSING")
	ErrMissingISBN      = errors.New("ERR_ISBN_MISSING")
	ErrMissingPublished = errors.New("ERR_PUBLISHED_MISSING")

	ErrInvalidID    = errors.New("ERR_INVALID_ID")
	ErrNotFoundBook = errors.New("ERR_BOOK_NOT_FOUND")

	ErrMissingUpdateData = errors.New("ERR_MISSING_UPDATE_DATA")
)
