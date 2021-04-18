// This is the server program that serves a course information microservice
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Course struct {
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Material struct {
	Id       string `json:"id"`
	Sequence int    `json:"sequence"`
	FileName string `json:"filename"`
}

var DB *sql.DB

// setupDB detects the environment and connects to the db accordingly
func setupDB() func() {
	var err error
	var dsn = "user:test@tcp(127.0.0.1:3306)/db"
	if os.Getenv("PLACE") == "DOCKER" {
		dsn = "user:test@tcp(db:3306)/db"
	}
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("db opened")
	}
	return func() {
		err = DB.Close()
		if err != nil {
			log.Println(err)
		} else {
			log.Println("db closed")
		}
	}
}

// setupLog sets up simple logging. Log file is given a different name everyday.
func setupLog() func() {
	t := time.Now()
	filename := fmt.Sprintf("log/Log%04d%02d%02d.txt", t.Year(), t.Month(), t.Day())
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("log opened")
	}
	wrt := io.MultiWriter(os.Stdout, file)
	log.SetOutput(wrt)
	return func() {
		err = file.Close()
		if err != nil {
			log.Println(err)
		} else {
			log.Println("log closed")
		}
	}
}

// rootHandler returns the index page
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET / - Return this page."))
	w.Write([]byte("GET /v1 - Return this page."))
	w.Write([]byte("GET /v1/courses – Return a list of all courses."))
	w.Write([]byte("POST /v1/courses – Add a new course."))
	w.Write([]byte("GET /v1/courses/{code} – Return the specified course."))
	w.Write([]byte("PUT /v1/courses/{code} – Update the specified course."))
	w.Write([]byte("DELETE /v1/courses/{code} – Delete the specified course."))
	w.Write([]byte("POST /v1/courses/{code}/materials – Add a new material to the specified course."))
	w.Write([]byte("GET /v1/courses/{code}/materials/metadata – Return metadata of all materials of the specified course."))
	w.Write([]byte("GET /v1/courses/{code}/materials/files – Download the specified material."))
	w.Write([]byte("PUT /v1/courses/{code}/materials/{id} – Update the specified material."))
	w.Write([]byte("DELETE /v1/courses/{code}/materials/{id} – Delete the specified material."))
	w.Write([]byte("GET /v1/courses/{code}/materials/{id}/metadata – Return the specified material."))
	w.Write([]byte("GET /v1/courses/{code}/materials/{id}/files – Download the specified material."))
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var v int
		// look up API in db
		q := "SELECT ISVALID FROM API WHERE API_ID=?"
		if err := DB.QueryRow(q, r.Header.Get("X-API-KEY")).Scan(&v); err != nil {
			log.Printf("QUERYDONE: %s, %s, %s\n", q, r.Header["X-API-KEY"], vars["code"])
			if fmt.Sprint(err) == "sql: no rows in result set" {
				log.Println(err)
				w.WriteHeader(http.StatusForbidden)
				return
			} else {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if v != 1 {
			log.Println("API access had been revoked")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// additional check for accessing specified course or material related, look up API_COURSE in db
		if _, ok := vars["code"]; ok == true {
			var a string
			q := "SELECT API_ID FROM API_COURSE WHERE API_ID=? AND COURSE_ID=?"
			if err := DB.QueryRow(q, r.Header.Get("X-API-KEY"), vars["code"]).Scan(&a); err != nil {
				log.Printf("QUERYDONE: %s, %s, %s\n", q, r.Header["X-API-KEY"], vars["code"])
				if fmt.Sprint(err) == "sql: no rows in result set" {
					log.Println(err)
					w.WriteHeader(http.StatusForbidden)
					return
				} else {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// main function starts this program
func main() {
	logCloser := setupLog()
	defer logCloser()
	dbCloser := setupDB()
	defer dbCloser()
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/v1", rootHandler).Methods("GET")
	r.HandleFunc("/v1/courses", ReturnAllCoursesHandler).Methods("GET")
	r.HandleFunc("/v1/courses", AddNewCourseHandler).Methods("POST")
	r.HandleFunc("/v1/courses/{code}", ReturnSpecificCourseHandler).Methods("GET")
	r.HandleFunc("/v1/courses/{code}", UpdateSpecificCourseHandler).Methods("PUT")
	r.HandleFunc("/v1/courses/{code}", RemoveSpecificCourseHandler).Methods("DELETE")
	r.HandleFunc("/v1/courses/{code}/materials", AddNewMaterialHandler).Methods("POST")
	r.HandleFunc("/v1/courses/{code}/materials/metadata", ReturnAllMaterialsHandler).Methods("GET")
	r.HandleFunc("/v1/courses/{code}/materials/files", DownloadAllMaterialsHandler).Methods("GET")
	r.HandleFunc("/v1/courses/{code}/materials/{id}", UpdateSpecificMaterialHandler).Methods("PUT")
	r.HandleFunc("/v1/courses/{code}/materials/{id}", RemoveSpecificMaterialHandler).Methods("DELETE")
	r.HandleFunc("/v1/courses/{code}/materials/{id}/metadata", ReturnSpecificMaterialHandler).Methods("GET")
	r.HandleFunc("/v1/courses/{code}/materials/{id}/files", DownloadSpecificMaterialHandler).Methods("GET")
	r.Use(authenticationMiddleware)
	http.Handle("/", r)
	if err := http.ListenAndServeTLS(":5000", "ssl/cert.pem", "ssl/key.pem", nil); err != nil {
		log.Fatalln(err)
	}
}
