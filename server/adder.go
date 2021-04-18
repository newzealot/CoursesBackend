package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// AddNewCourseHandler add a new course.
func AddNewCourseHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := receiveJSON(r)
	log.Printf("RECEIVED: %#v\n", c)
	// insert course into db
	q := "INSERT INTO COURSE(COURSE_ID,TITLE,DESCRIPTION) VALUES (?,?,?)"
	if err := ExecQuery(w, q, c.Code, c.Title, c.Description); err != nil {
		return
	}
	// grant course access to API Key
	q = "INSERT INTO API_COURSE(API_ID,COURSE_ID) VALUES (?,?)"
	if err := ExecQuery(w, q, r.Header.Get("X-API-KEY"), c.Code); err != nil {
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// AddNewMaterialHandler adds a new material to the specified course.
func AddNewMaterialHandler(w http.ResponseWriter, r *http.Request) {
	var m Material
	vars := mux.Vars(r)
	err := r.ParseMultipartForm(128 << 20)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("RECEIVED: %#v\n", r.FormValue("json"))
	}
	// receiving material metadata
	if err := json.Unmarshal([]byte(r.FormValue("json")), &m); err != nil {
		log.Println(err)
	}
	// insert into db
	q := "INSERT INTO MATERIAL(MATERIAL_ID,SEQUENCE,FILENAME,COURSE_ID) VALUES (?,?,?,?)"
	if err := ExecQuery(w, q, uuid.NewString(), m.Sequence, m.FileName, vars["code"]); err != nil {
		return
	}
	if err := formFileToDisk(w, r, vars["code"], m); err != nil {
		return
	}
	// rezip folder to allow retrievals of all materials of specific course
	zipUpMaterials(vars["code"])
	w.WriteHeader(http.StatusAccepted)
}
