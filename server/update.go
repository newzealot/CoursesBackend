package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// UpdateSpecificCourseHandler updates the specified course.
// You cannot change the COURSE_ID.
// If you really need to, delete and add new course again.
func UpdateSpecificCourseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, _ := receiveJSON(r)
	log.Printf("RECEIVED: %#v\n", c)
	// update course in db
	q := "UPDATE COURSE SET TITLE=?, DESCRIPTION=? WHERE COURSE_ID=?"
	if err := ExecQuery(w, q, c.Title, c.Description, vars["code"]); err != nil {
		return
	}
	log.Printf("QUERYDONE: %s, %s\n", q, c.Code)
	// sends 200 OK by default
}

// UpdateSpecificMaterialHandler updates the specified material.
// You cannot change the MATERIAL_ID.
// If you really need to, delete and add new material again.
func UpdateSpecificMaterialHandler(w http.ResponseWriter, r *http.Request) {
	var m1, m2 Material
	vars := mux.Vars(r)
	if err := r.ParseMultipartForm(128 << 20); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("RECEIVED: %#v\n", r.FormValue("json"))
	// receive incoming file metadata
	if err := json.Unmarshal([]byte(r.FormValue("json")), &m1); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// getting sequence and filename of existing material
	// this is to process the removal of old files
	q := "SELECT SEQUENCE, FILENAME FROM MATERIAL WHERE MATERIAL_ID=?"
	if err := DB.QueryRow(q, vars["id"]).Scan(&m2.Sequence, &m2.FileName); err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("QUERYDONE: %s, %s\n", q, m1.Id)

	// update material in db
	q = "UPDATE MATERIAL SET SEQUENCE=?, FILENAME=? WHERE MATERIAL_ID=?"
	if err := ExecQuery(w, q, m1.Sequence, m1.FileName, vars["id"]); err != nil {
		return
	}
	// remove old files
	loc := fmt.Sprintf("download/%s/%d/%s", vars["code"], m2.Sequence, m2.FileName)
	if err := removeLoc(loc); err != nil {
		w.WriteHeader(http.StatusInsufficientStorage)
		return
	}
	if err := formFileToDisk(w, r, vars["code"], m1); err != nil {
		return
	}
	// zip up code folder
	zipUpMaterials(vars["code"])
	// sends 200 OK by default
}
