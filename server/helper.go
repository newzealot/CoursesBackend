package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// receiveJSON process receiving incoming json data
func receiveJSON(r *http.Request) (Course, Material) {
	var c Course
	var m Material
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
	}
	log.Printf("RECEIVED: %s\n", string(body))
	if strings.Contains(string(body), `"title":`) {
		err = json.Unmarshal(body, &c)
		if err != nil {
			log.Println(err)
		}
	} else if strings.Contains(string(body), `"sequence":`) {
		err = json.Unmarshal(body, &m)
		if err != nil {
			log.Println(err)
		}
	}
	return c, m
}

// formFileToDisk saves form file to disk
func formFileToDisk(w http.ResponseWriter, r *http.Request, v string, m Material) error {
	// receiving material file
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	loc := fmt.Sprintf("download/%s/%d", v, m.Sequence)
	err = os.MkdirAll(loc, 0755)
	if err != nil {
		log.Println(err)
		return err
	}
	loc = fmt.Sprintf("%s/%s", loc, handler.Filename)
	f, err := os.OpenFile(loc, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInsufficientStorage)
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	_, err = io.Copy(f, file)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInsufficientStorage)
		return err
	}
	return nil
}

// zipUpMaterials zips up the COURSE_ID folder
func zipUpMaterials(code string) {
	// create a new zip file
	loc := fmt.Sprintf("download/%s", code)
	file, err := os.OpenFile(loc+".zip", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		file.Close()
		log.Printf("FILECREATED: %s.zip\n", loc)
	}()
	w := zip.NewWriter(file)
	defer w.Close()
	// walks the filepath to gather all valid files
	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// opens valid file
		file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			return err
		}
		defer file.Close()
		// names file in zip
		filenameInZip := strings.SplitAfterN(path, "/", 4)
		f, err := w.Create(filenameInZip[len(filenameInZip)-1])
		if err != nil {
			return err
		}
		// copy valid file to .zip file
		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
		return nil
	}
	err = filepath.Walk(loc, walker)
	if err != nil {
		log.Println(err)
	}
}

// removeLoc removes a file
func removeLoc(loc string) error {
	if err := os.RemoveAll(loc); err != nil {
		log.Printf("ERROR: %s\n", err)
		return err
	}
	log.Printf("FILEDELETED: %s\n", loc)
	return nil
}
