package models

import (
	"database/sql"
	"fmt"
	"log"
)

type Quote struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
	Id     int    `json:"id"`
}

type QuoteModel struct {
	DB *sql.DB
}

func (qm *QuoteModel) Insert(quote string, author string) (int, error) {
	var id int
	queryStatment := `INSERT INTO quotes (quote, author) VALUES ($1, $2) RETURNING ID;`
	err := qm.DB.QueryRow(queryStatment, quote, author).Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	return int(id), nil
}

func (qm *QuoteModel) Get() ([]Quote, error) {
	queryStatment := `SELECT * FROM quotes;`

	var quotes []Quote

	rows, err := qm.DB.Query(queryStatment)
	if err != nil {
		fmt.Println("db query")
		log.Println(err.Error())
		return []Quote{}, err
	}

	for rows.Next() {
		var dbRow Quote
		err := rows.Scan(&dbRow.Id, &dbRow.Quote, &dbRow.Author)
		if err != nil {
			log.Print(err.Error(), "rows next")
			return []Quote{}, err
		}

		quotes = append(quotes, dbRow)
	}

	err2 := rows.Err()
	if err2 != nil {
		log.Print(err2.Error(), "rows err")
		return []Quote{}, err2
	}

	return quotes, nil
}

func (qm *QuoteModel) GetQuote(id int) (Quote, error) {
	queryStatment := `SELECT id,quote,author FROM quotes WHERE id = $1;`

	var quote Quote

	err := qm.DB.QueryRow(queryStatment, id).Scan(&quote.Id, &quote.Quote, &quote.Author)
	if err != nil {
		log.Println(err.Error())
		return Quote{}, err
	}
	return quote, nil
}
