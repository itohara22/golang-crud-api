package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	res.Write([]byte("Hello"))
}

type Quote struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
	Id     int    `json:"id"`
}

func (app *application) getQuotes(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	queryId := req.URL.Query().Get("id")

	res.Header().Set("Content-Type", "application/json")

	if queryId != "" {
		id, err := strconv.Atoi(queryId)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		quote, err := app.quotes.GetQuote(int(id))
		if err != nil {
			http.NotFound(res, req)
			return
		}
		err2 := json.NewEncoder(res).Encode(quote)
		if err2 != nil {
			http.Error(res, "something went wrong", http.StatusInternalServerError)
			return
		}

	} else {
		quotes, err := app.quotes.Get()
		if err != nil {
			http.Error(res, "something went wrong", http.StatusInternalServerError)
			return
		}

		err2 := json.NewEncoder(res).Encode(quotes)
		if err2 != nil {
			http.Error(res, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
}

func (app *application) createQuote(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var newQuote Quote

	err := json.NewDecoder(req.Body).Decode(&newQuote)
	if err != nil {
		http.Error(res, "something went wrong", http.StatusInternalServerError)
		return
	}

	id, err := app.quotes.Insert(newQuote.Quote, newQuote.Author)
	if err != nil {
		http.Error(res, "something went wrong", http.StatusInternalServerError)
		return
	}

	newQuote.Id = id

	json.NewEncoder(res).Encode(newQuote)
}

func (app *application) handleQuotes(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		app.getQuotes(res, req)
	case http.MethodPost:
		app.createQuote(res, req)
	default:
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
