package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

// DownloadAllMaterialsHandler downloads files of all materials of the specified course.
func DownloadAllMaterialsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loc := fmt.Sprintf("download/%s", vars["code"])
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
	} else {
		loc := fmt.Sprintf("download/%s.zip", vars["code"])
		log.Printf("ACTION: Serving file %s\n", loc)
		w.Header().Add("Filename", vars["code"]+".zip")
		http.ServeFile(w, r, loc)
	}
}

// DownloadSpecificMaterialHandler downloads the specified material.
func DownloadSpecificMaterialHandler(w http.ResponseWriter, r *http.Request) {
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
		loc := fmt.Sprintf("download/%s/%d/%s", vars["code"], m.Sequence, m.FileName)
		log.Printf("ACTION: Serving file %s\n", loc)
		w.Header().Add("Filename", m.FileName)
		http.ServeFile(w, r, loc)
	}
}
