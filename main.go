package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ilhamsyahids/bookshelf-template/rest"
	"github.com/ilhamsyahids/bookshelf-template/storage"
	"github.com/jmoiron/sqlx"
)

const addr = ":8080"

const envSQLDSN = "SQL_DSN"

func main() {
	sqlDSN := os.Getenv(envSQLDSN)
	sqlClient, err := sqlx.Connect("mysql", sqlDSN)
	if err != nil {
		log.Fatalf("unable to connect to database due: %v", err)
	}
	defer sqlClient.Close()

	// init storage
	storage := storage.NewStorage(sqlClient)

	// init API
	api, err := rest.NewAPI(rest.APIConfig{BookStorage: *storage})
	if err != nil {
		log.Fatalf("unable to initialize api due: %v", err)
	}

	// init server
	server := &http.Server{
		Addr:    addr,
		Handler: api.GetHandler(),
	}
	// run server
	log.Printf("server is listening on %v", addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("unable to run server due: %v", err)
	}
}
