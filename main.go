package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"
)

type Rate struct {
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Mid      float64 `json:"mid"`
}

type ExchangeTable struct {
	Table         string `json:"table"`
	No            string `json:"no"`
	EffectiveDate string `json:"effectiveDate"`
	Rates         []Rate `json:"rates"`
}

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func main() {

	client := &http.Client{
		Timeout: time.Second * 10, // Timeout each requests
	}
	url := "https://api.nbp.pl/api/exchangerates/tables/A/"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Err ", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Err: ", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Err: ", err)
		return
	}
	fmt.Println(string(body))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			panic(err)
		}
	})

	http.ListenAndServe(":8080", nil)
}
