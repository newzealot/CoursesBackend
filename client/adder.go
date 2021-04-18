package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// AddNewCourseHandler adds a new course.
func AddNewCourseHandler() {
	fmt.Println("1. Add a new course.")
	c, _, err := getInput(CourseCode, CourseTitle, CourseDescription)
	if err != nil {
		log.Println("error getting input. restarting")
		return
	}
	requestBody, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Println(err)
	}
	sendRequest(http.MethodPost, CourseURL, ContentTypeJSON, requestBody)
}

// AddNewMaterialHandler adds a new material to the specified course.
// The material to be uploaded has to be stored in the upload folder.
func AddNewMaterialHandler() {
	fmt.Println("2. Add a new material to the specified course")
	c, m, err := getInput(CourseCode, MaterialSequence, MaterialFileName)
	if err != nil {
		log.Println("error getting input. restarting")
		return
	}
	s := fmt.Sprintf("%s/%s/materials", CourseURL, c.Code)
	j, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		log.Println(err)
	}
	ct, b, err := createMultiPart(m, j)
	if err != nil {
		return
	}
	sendRequest(http.MethodPost, s, ct, b)
}
