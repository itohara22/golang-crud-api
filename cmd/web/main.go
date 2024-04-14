package main

import (
	"context"
	"database/sql"
	"flag"
	"idontKnowWhatIamDoing/internal/models"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	quotes *models.QuoteModel
}

func main() {
	cliPortFlag := flag.String("port", ":6969", "HTTP network address")
	flag.Parse()

	db, err2 := connectToDb()
	if err2 != nil {
		log.Fatal("cant connect to db")
		return
	}

	defer db.Close()

	// we use this struct th make QuoteModel qury methods available to our handlers
	app := &application{
		quotes: &models.QuoteModel{DB: db},
	}

	router := http.NewServeMux()

	router.HandleFunc("/", app.home)
	router.HandleFunc("/quotes", app.handleQuotes)

	// we define a  server struct
	server := &http.Server{
		Addr:    *cliPortFlag,
		Handler: router,
	}

	log.Printf("Server started on %s", *cliPortFlag)
	// err := http.ListenAndServe(*cliPortFlag, router)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func connectToDb() (*sql.DB, error) {
	connStr := "postgres://postgres:postgres@localhost:8080/mydatabase"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	err2 := db.PingContext(context.Background())
	if err2 != nil {
		log.Fatal(err2.Error())
		return nil, err2
	}

	return db, nil
}
