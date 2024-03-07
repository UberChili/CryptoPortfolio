package main

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var store = sessions.NewCookieStore([]byte("super-secret"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/index.html"))
	session, _ := store.Get(r, "session")
	_, ok := session.Values["user_id"]
	if !ok {
		http.Redirect(w, r, "/login/", http.StatusFound)
		return
	}

	// tmpl.Execute(w, nil)
	fmt.Println(session.Values["user_id"])
	tmpl.Execute(w, session.Values)
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
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/login.html"))

	if r.Method != "POST" { // User reached route via GET
		tmpl.Execute(w, nil)
		return
	}
	// User reached route via POST (as by submitting a form via PSOT)
	// Ensure user provided username and password
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing!")
		return
	}
	user := r.PostFormValue("username")
	pass := r.PostFormValue("password")
	if user == "" || pass == "" {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	// Check if user exists
	// stmt := "SELECT id FROM users WHERE username = ?"
	stmt := "SELECT id, username, hash FROM users WHERE username = ?"
	row := db.QueryRow(stmt, user)
	var id, username, hash string
	err_scan := row.Scan(&id, &username, &hash)
	if err_scan == sql.ErrNoRows {
		fmt.Println("User does not exist")
		err_message := "User does not exist"
		Message := map[string]string{"Message": err_message}
		tmpl.Execute(w, Message)
		return
	}
	// Check if password is correct.
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		fmt.Println("Password is incorrect")
		err_message := "Incorrect username or password"
		Message := map[string]string{"Message": err_message}
		tmpl.Execute(w, Message)
		return
	}
	// Logging in was succes, save session and redirect to index
	session, _ := store.Get(r, "session")
	session.Values["user_id"] = id
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

// Function to handle the register page
func registerHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/register.html"))
	if r.Method != "POST" { // User reached route via POST (as by submitting a form via POST)
		tmpl.Execute(w, nil)
		return
	}

	// Parse form values
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}
	username := r.PostFormValue("username")
	pass := r.PostFormValue("password")
	confirmation := r.PostFormValue("confirmation")

	// Check for correct value criteria
	var err_message string

	if username == "" || pass == "" || confirmation == "" {
		err_message = "Username, password and password confirmation are required"
	} else if pass != confirmation {
		err_message = "Password and password confirmation are not the same"
	}
	if err_message != "" {
		var Message = map[string]string{"Message": err_message}
		fmt.Println(err_message)
		tmpl.Execute(w, Message)
		return
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
	// Redirect to login page or display success message
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

	// Handle functions
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/quote/", quoteHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/register/", registerHandler)

	log.Println("App running on 8000...")
	// log.Fatal(http.ListenAndServe(":8000", nil))
	log.Fatal(http.ListenAndServe("localhost:8000", context.ClearHandler(http.DefaultServeMux)))
}
