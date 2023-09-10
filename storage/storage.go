package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Book struct {
	ID        string `json:"id"`
	ISBN      string `json:"isbn"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Published string `json:"published"`
}

type Storage struct {
	booksData []Book
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Load(filename string) error {
	rawBook, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("unable to read file due: %v", err)
	}

	err = json.Unmarshal(rawBook, &s.booksData)
	if err != nil {
		return fmt.Errorf("unable to init books data due: %v", err)
	}

	return nil
}

func (s *Storage) GetBooks() ([]Book, error) {
	return s.booksData, nil
}

func (s *Storage) AddBook(book Book) error {
	s.booksData = append(s.booksData, book)
	return nil
}

func (s *Storage) GetBookByID(id string) (Book, error) {
	for _, book := range s.booksData {
		if book.ID == id {
			return book, nil
		}
	}
	return Book{}, fmt.Errorf("book not found")
}

func (s *Storage) UpdateBookByID(id string, book Book) error {
	for i, b := range s.booksData {
		if b.ID == id {
			// Update the book at the given index
			s.booksData[i] = book
			return nil
		}
	}
	return fmt.Errorf("book not found")
}

func (s *Storage) DeleteBookByID(id string) error {
	for i, book := range s.booksData {
		if book.ID == id {
			s.booksData = append(s.booksData[:i], s.booksData[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("book not found")
}
