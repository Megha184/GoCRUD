package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Movie struct {
	ID    string `json:"id"`
	Isbn  string `json:"isbn"`
	Title string `json:"title"`
}

const (
	host      = "localhost"
	port      = 5432
	user      = "postgres"
	password  = "abc123"
	dbname    = "moviedb"
	tablename = "movietable"
)

var db *sql.DB

func main() {
	/* Sprintf formats according to a format specifier and returns the string */
	psqlconn := fmt.Sprintf("host=%s port =%d user= %s password =%s dbname =%s sslmode = disable", host, port, user, password, dbname)
	//open the database so we can perform the task
	db, _ = sql.Open(user, psqlconn)
	fmt.Println("Connected!")
	/** Started frome here **/
	// New Router
	r := mux.NewRouter()
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	// Print starting at port 8000
	fmt.Printf("Starting server at port 8000\n")
	//	start the server
	log.Fatal(http.ListenAndServe(":8000", r))
	// Closing the connection to avoid any kind of memory overflow or run time error panic state
	defer db.Close()

}

// Update the Movie and also print the movie for me to see it is updated
func updateMovie(w http.ResponseWriter, r *http.Request) {
	//Response header we are giving response in json format
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var mov Movie
	// r.body is all the json data we added in postman
	// but we need to convert it into movie type
	err := json.NewDecoder(r.Body).Decode(&mov)
	if err != nil {
		fmt.Print("unable to decode, maybe json not in correct format")
	}
	exe := `update movietable set isbn = $1, title = $2 where id = $3`
	_, err = db.Query(exe, mov.Isbn, mov.Title, params["id"])
	if err != nil {
		fmt.Println("Not updated error Occured")
	} else {
		exe = `select * from movietable where id=$1`
		row, _ := db.Query(exe, params["id"])
		for row.Next() {
			row.Scan(&mov.ID, &mov.Isbn, &mov.Title)
			mov.ID = params["id"]
		}
		defer row.Close()
		json.NewEncoder(w).Encode(mov)
	}
	// Why this is not working?
	//db.Query("update movietable set isbn = ?, title = ? where id = ?", params["Isbn"], params["title"], params["ID"])

}

// Display All the details from the database
func getMovies(w http.ResponseWriter, r *http.Request) {
	var movies []Movie
	movies = nil
	w.Header().Set("Content-Type", "application/json")
	row, _ := db.Query("Select * from movietable")
	for row.Next() {
		var mov Movie
		row.Scan(&mov.ID, &mov.Isbn, &mov.Title)
		movies = append(movies, mov)
	}
	defer row.Close()
	json.NewEncoder(w).Encode(movies)
	movies = nil
}

// Create A particular Movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var mov Movie
	_ = json.NewDecoder(r.Body).Decode(&mov)
	insertDynStmt := `insert into "movietable"("id", "isbn","title") values($1, $2, $3)`
	_, _ = db.Exec(insertDynStmt, mov.ID, mov.Isbn, mov.Title)
}

// Delete a movie using id
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	fmt.Println(params, " is deleted")
	exe := `delete from movietable where id = $1`
	db.Exec(exe, params["id"])
}

// Get a particular movie using a particular id
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	st := `select * from movietable where id = $1`
	// row doesn't work as row just has the address so we need to use row.Next() to get the actual value
	// Scan is used to scan that row in a partcular fomrat of movie
	row, _ := db.Query(st, params["id"])
	var mov Movie
	for row.Next() {

		row.Scan(&mov.ID, &mov.Isbn, &mov.Title)
	}
	defer row.Close()
	json.NewEncoder(w).Encode(mov)
}
