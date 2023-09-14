package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var BookNotFound = fmt.Errorf("book not found")

type Book struct {
	ID        string `json:"id"`
	ISBN      string `json:"isbn"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Published string `json:"published"`
}

type Storage struct {
	sqlClient *sqlx.DB
}

func NewStorage(sqlClient *sqlx.DB) *Storage {
	return &Storage{sqlClient: sqlClient}
}

func (s *Storage) GetBooks(query string, page, limit int) ([]Book, error) {
	// get books from database
	books := []Book{}
	sqlQuery := `SELECT * FROM books WHERE title LIKE ? LIMIT ? OFFSET ?`
	err := s.sqlClient.Select(&books, sqlQuery, "%"+query+"%", limit, (page - 1))
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *Storage) CreateBook(book Book) (*Book, error) {
	// insert book to database
	sqlQuery := `INSERT INTO books (isbn, title, author, published) VALUES (?, ?, ?, ?)`
	res, err := s.sqlClient.Exec(sqlQuery, book.ISBN, book.Title, book.Author, book.Published)
	if err != nil {
		return nil, err
	}

	// get last inserted id
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	// set id to book
	book.ID = fmt.Sprintf("%d", id)

	return &book, nil
}

func (s *Storage) GetBookByID(id string) (*Book, error) {
	// get book from database
	book := Book{}
	sqlQuery := `SELECT * FROM books WHERE id = ?`
	err := s.sqlClient.Get(&book, sqlQuery, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, BookNotFound
		}
		return nil, err
	}
	return &book, nil
}

func (s *Storage) UpdateBook(id string, book Book) (*Book, error) {
	// update book from database
	sqlQuery := `UPDATE books SET`
	args := []interface{}{}
	if book.ISBN != "" {
		sqlQuery += ` isbn = ?,`
		args = append(args, book.ISBN)
	}
	if book.Title != "" {
		sqlQuery += ` title = ?,`
		args = append(args, book.Title)
	}
	if book.Author != "" {
		sqlQuery += ` author = ?,`
		args = append(args, book.Author)
	}
	if book.Published != "" {
		sqlQuery += ` published = ?,`
		args = append(args, book.Published)
	}

	// remove last comma
	sqlQuery = sqlQuery[:len(sqlQuery)-1]

	// add where clause
	sqlQuery += ` WHERE id = ?`
	args = append(args, id)

	_, err := s.sqlClient.Exec(sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	// get updated book
	return s.GetBookByID(id)
}

func (s *Storage) DeleteBook(id string) error {
	// delete book from database
	sqlQuery := `DELETE FROM books WHERE id = ?`
	_, err := s.sqlClient.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}
