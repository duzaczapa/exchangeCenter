package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Rate struct {
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Mid      float64 `json:"mid"`
}

type ExchangeTable struct {
	Table string `json:"table"`
	Rates []Rate `json:"rates"`
}

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func fetchRates() (*ExchangeTable, error) {
	url := "https://api.nbp.pl/api/exchangerates/tables/A/"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Err ", err)
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Err: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: %d", resp.StatusCode)
	}

	var tables []ExchangeTable
	if err := json.NewDecoder(resp.Body).Decode(&tables); err != nil {
		return nil, err
	}

	if len(tables) == 0 {
		return nil, fmt.Errorf("There are no rates here")
	}

	return &tables[0], nil
}

func main() {

	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		table, _ := fetchRates()

		if err := tmpl.Execute(w, table); err != nil {
			log.Printf("Err: %v", err)
		}
	})

	http.ListenAndServe(":8080", nil)
}
