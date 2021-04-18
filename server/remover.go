package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// RemoveSpecificCourseHandler deletes the specified course.
func RemoveSpecificCourseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	q := "DELETE FROM COURSE WHERE COURSE_ID = ?"
	if err := ExecQuery(w, q, vars["code"]); err != nil {
		return
	}
	// remove specified course material
	loc := fmt.Sprintf("download/%s", vars["code"])
	if err := removeLoc(loc); err != nil {
		w.WriteHeader(http.StatusInsufficientStorage)
		return
	}
	// remove zip file
	loc = fmt.Sprintf("download/%s.zip", vars["code"])
	if err := removeLoc(loc); err != nil {
		w.WriteHeader(http.StatusInsufficientStorage)
		return
	}
	// sends 200 OK by default
}

// RemoveSpecificMaterialHandler deletes the specified material.
func RemoveSpecificMaterialHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var m Material
	var c Course
	// taking note of old material metadata so that can delete file
	q2 := "Select * FROM MATERIAL WHERE MATERIAL_ID = ?"
	log.Printf("QUERY: %s, %s\n", q2, vars["id"])
	err := DB.QueryRow(q2, vars["id"]).Scan(&m.Id, &m.Sequence, &m.FileName, &c.Code)
	log.Println("QUERY closed")
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	// delete material files from db
	q := "DELETE FROM MATERIAL WHERE MATERIAL_ID = ?"
	if err := ExecQuery(w, q, vars["id"]); err != nil {
		return
	}
	// remove specified course material
	loc := fmt.Sprintf("download/%s/%d", vars["code"], m.Sequence)
	if err := removeLoc(loc); err != nil {
		w.WriteHeader(http.StatusInsufficientStorage)
		return
	}
	// remove zip file
	loc = fmt.Sprintf("download/%s.zip", vars["code"])
	if err := removeLoc(loc); err != nil {
		w.WriteHeader(http.StatusInsufficientStorage)
		return
	}
	zipUpMaterials(vars["code"])
	// sends 200 OK by default
}
