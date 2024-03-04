package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type QuoteResponse struct {
	Data map[string]QuoteData `json:"data"`
}

type Quote struct {
	USD QuoteDetails `json:"USD"`
}

type QuoteDetails struct {
	Price     float64 `json:"price"`
	MarketCap float64 `json:"market_cap"`
}

type QuoteData struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Symbol string `json:"symbol"`
	Quote  Quote  `json:"quote"`
}

type Simple_quote struct {
	Name   string
	Slug   string
	Symbol string
	Price  float64
}

const CMC_API_KEY = "8dc4e9af-b945-4a13-b5d7-a315ed95fce9"

func get_quote(slug string) Simple_quote {
	slug = strings.ToLower(slug)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("convert", "USD")
	q.Add("slug", slug)
	// q.Add("slug", "ethereum")
	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", CMC_API_KEY)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	} else if resp.StatusCode == 400 {
		fmt.Println("Coin not found")
		return Simple_quote{}
	}
	// fmt.Println(resp.Status)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in response")
		os.Exit(1)
	}

	var quote QuoteResponse
	errq := json.Unmarshal(respBody, &quote)
	if errq != nil {
		fmt.Println("Error: Couldn't unmarshal")
		os.Exit(1)
	}

	var simple_quote Simple_quote

	for _, v := range quote.Data {
		simple_quote.Name = v.Name
		simple_quote.Slug = v.Slug
		simple_quote.Symbol = v.Symbol
		// simple_quote.Price = v.Quote.USD.Price
		simple_quote.Price = float64(int(v.Quote.USD.Price*100)) / 100
	}
	// fmt.Println(simple_quote.Name)
	// fmt.Println(simple_quote.Slug)
	// fmt.Println(simple_quote.Symbol)
	// fmt.Println(simple_quote.Price)
	return simple_quote
}
