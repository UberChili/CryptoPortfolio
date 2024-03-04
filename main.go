package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	// "golang.org/x/crypto/bcrypt"
)

type Session struct {
	Username string
	id       string
	Logged   bool
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/index.html"))
	tmpl.Execute(w, nil)
}

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", indexHandler)
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/index.html"))
	// 	tmpl.Execute(w, nil)
	// })

	/* Function to Quote cryptocurrencies
	Deals with both get and post requests */
	http.HandleFunc("/quote/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Post request
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form.", http.StatusBadRequest)
				return
			}
			// Ensure user provides slug
			slug := r.PostFormValue("slug")
			if slug == "" {
				http.Error(w, "Must provide slug.", http.StatusBadRequest)
				return
			}
			data := map[string]Simple_quote{
				"Results": get_quote(slug),
			}
			for _, v := range data {
				fmt.Println(v.Name, v.Symbol, v.Slug, v.Price)
			}
			// tmpl := template.Must(template.ParseFiles("./templates/quoted.html"))
			tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/quoted.html"))
			tmpl.Execute(w, data)
		} else {
			// Get request
			// tmpl := template.Must(template.ParseFiles("./templates/quote.html"))
			tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/quote.html"))
			tmpl.Execute(w, nil)
		}
	})

	http.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" { // User reached route via POST (as by submitting a form via POST)
			// Ensure user provided username and password
			err := r.ParseForm()
			if err != nil {
				fmt.Println("Error parsing!")
				return
			}
			name := r.PostFormValue("username")
			pass := r.PostFormValue("password")
			if name == "" || pass == "" {
				fmt.Println("Must provide username and password.")
				return
			}
			fmt.Println(name, pass)
			fmt.Println("Got to make it look like I'm doing work or else my dad will have me doing some stupid shit")
		} else { // User reached route via GET (as by clicking a link or via redirect)
			tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/login.html"))
			tmpl.Execute(w, nil)
		}
	})

	// http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl := template.Must(template.ParseFiles("./templates/fragments/results.html"))
	// 	data := map[string][]Stock{
	// 		"Results": SearchTicker(r.URL.Query().Get("key")),
	// 	}
	// 	tmpl.Execute(w, data)
	// })

	// http.HandleFunc("/stock/", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case "POST":
	// 		ticker := r.PostFormValue("ticker")
	// 		stk := SearchTicker(ticker)[0]
	// 		val := GetDailyValues(ticker)
	// 		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	// 		tmpl.ExecuteTemplate(w, "stock-element",
	// 			Stock{Ticker: stk.Ticker, Name: stk.Name, Price: val.Open})
	// 	}
	// })

	log.Println("App running on 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
