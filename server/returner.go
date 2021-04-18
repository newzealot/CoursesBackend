package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// ReturnAllCoursesHandler returns a list of all courses.
func ReturnAllCoursesHandler(w http.ResponseWriter, r *http.Request) {
	var a []Course
	q := "Select * FROM COURSE"
	if err := QueryAll(w, &a, q); err != nil {
		return
	}
	requestBody, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		log.Println(err)
	}
	log.Println(fmt.Sprintf("RESULT: %#v", string(requestBody)))
	if _, err := w.Write(requestBody); err != nil {
		log.Println(err)
	}
}

// ReturnSpecificCourseHandler returns the specified course.
func ReturnSpecificCourseHandler(w http.ResponseWriter, r *http.Request) {
	c := Course{}
	vars := mux.Vars(r)
	q := "Select * FROM COURSE WHERE COURSE_ID = ?"
	log.Printf("QUERY: %s, %s\n", q, vars["code"])
	err := DB.QueryRow(q, vars["code"]).Scan(&c.Code, &c.Title, &c.Description)
	log.Println("QUERY closed")
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		requestBody, err := json.MarshalIndent(c, "", "    ")
		if err != nil {
			log.Println(err)
		}
		s := fmt.Sprintf("RESULT: %#v", string(requestBody))
		log.Println(s)
		if _, err := w.Write(requestBody); err != nil {
			log.Println(err)
		}
	}
}

// ReturnAllMaterialsHandler returns metadata of all materials of the specified course.
func ReturnAllMaterialsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var a []Material
	q := "Select MATERIAL_ID, SEQUENCE, FILENAME FROM MATERIAL WHERE COURSE_ID = ?"
	if err := QueryAll(w, &a, q, vars["code"]); err != nil {
		return
	}
	requestBody, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		log.Println(err)
	}
	log.Println(fmt.Sprintf("RESULT: %#v", string(requestBody)))
	if _, err := w.Write(requestBody); err != nil {
		log.Println(err)
	}
}

// ReturnSpecificMaterialHandler returns the specified material.
func ReturnSpecificMaterialHandler(w http.ResponseWriter, r *http.Request) {
	m := Material{}
	c := Course{}
	vars := mux.Vars(r)
	q := "Select * FROM MATERIAL WHERE MATERIAL_ID = ?"
	log.Printf("QUERY: %s, %s\n", q, vars["id"])
	err := DB.QueryRow(q, vars["id"]).Scan(&m.Id, &m.Sequence, &m.FileName, &c.Code)
	log.Println("QUERY closed")
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		requestBody, err := json.MarshalIndent(m, "", "    ")
		if err != nil {
			log.Println(err)
		}
		s := fmt.Sprintf("RESULT: %#v", string(requestBody))
		log.Println(s)
		if _, err := w.Write(requestBody); err != nil {
			log.Println(err)
		}
	}
}
