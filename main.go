package main

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/index.html"))
	tmpl.Execute(w, nil)
}

// Function to Quote cryptocurrencies
// Deals with both get and post requests
func quoteHandler(w http.ResponseWriter, r *http.Request) {
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
}

// Function to handle the login page
func loginHandler(w http.ResponseWriter, r *http.Request) {
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
		// fmt.Println(name, pass)
	} else { // User reached route via GET (as by clicking a link or via redirect)
		tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/login.html"))
		tmpl.Execute(w, nil)
	}
}

// Function to handle the register page
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" { // User reached route via POST (as by submitting a form via POST)
		// Ensure user provided username and password
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Error parsing!")
			return
		}
		// Get values
		username := r.PostFormValue("username")
		pass := r.PostFormValue("password")
		confirmation := r.PostFormValue("confirmation")
		fmt.Println(username, pass, confirmation)

		// Check for correct value criteria
		var err_message string

		if username == "" {
			fmt.Println("Username field empty")
			err_message = "Must provide a username"
		} else if pass == "" || confirmation == "" {
			fmt.Println("Password or confirmation are empty")
			err_message = "Must provide a password and a password confirmation"
		} else if pass != confirmation {
			fmt.Println("Different password and confirmation")
			err_message = "Password and password confirmation are not the same"
		}
		if err_message != "" {
			var Message = map[string]string{"Message": err_message}
			tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/register.html"))
			tmpl.Execute(w, Message)
		}
		// Check if user already exists
		stmt := "SELECT id FROM users WHERE username = ?"
		row := db.QueryRow(stmt, username)
		var id string
		err_scan := row.Scan(&id)
		if err_scan != sql.ErrNoRows {
			fmt.Println("User already exists")
			err_message = "User already exists"
			var Message = map[string]string{"Message": err_message}
			tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/register.html"))
			tmpl.Execute(w, Message)
			return
		}
		// Create hash from password
		var hash []byte
		hash, err = bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("bcrypt err: ", err)
			err_message = err.Error()
			var Message = map[string]string{"Message": err_message}
			tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/register.html"))
			tmpl.Execute(w, Message)
			return
		}
		// Insert new user into table
		var insertStmt *sql.Stmt
		insertStmt, err = db.Prepare("INSERT INTO users (username, hash) VALUES (?, ?)")
		if err != nil {
			fmt.Println("Error preparing insert statement")
			return
		}
		defer insertStmt.Close()
		var result sql.Result
		result, err = insertStmt.Exec(username, hash)
		rowsAff, _ := result.RowsAffected()
		lastIns, _ := result.LastInsertId()
		fmt.Println("RowsAff: ", rowsAff)
		fmt.Println("lastIns: ", lastIns)
		fmt.Println("err: ", err)
		if err != nil {
			fmt.Println("Error inserting new user to table")
			return
		}

	} else { // User reached route via GET (as by clicking a link or via redirect)
		tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/register.html"))
		tmpl.Execute(w, nil)
	}
}

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	var err error
	db, err = sql.Open("sqlite3", "finance.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// rows, err := db.Query("SELECT * FROM users")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for rows.Next() {
	// 	var id int
	// 	var username string
	// 	var hash string
	// 	var cash float64
	//
	// 	err = rows.Scan(&id, &username, &hash, &cash)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("%d, %s, %s, %f\n", id, username, hash, cash)
	// }

	// Handle functions
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/quote/", quoteHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/register/", registerHandler)

	log.Println("App running on 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
