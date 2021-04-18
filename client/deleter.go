package main

import (
	"fmt"
	"log"
	"net/http"
)

// RemoveSpecificCourseHandler deletes the specified course.
func RemoveSpecificCourseHandler() {
	fmt.Println("Delete a specified course.")
	c, _, err := getInput(CourseCode)
	if err != nil {
		log.Println("error getting input. restarting")
		return
	}
	s := fmt.Sprintf("%s/%s", CourseURL, c.Code)
	sendRequest(http.MethodDelete, s, "", nil)
}

// RemoveSpecificMaterialHandler deletes the specified material.
func RemoveSpecificMaterialHandler() {
	fmt.Println("6. Delete a specified material.")
	c, m, err := getInput(CourseCode, MaterialId)
	if err != nil {
		log.Println("error getting input. restarting")
		return
	}
	s := fmt.Sprintf("%s/%s/materials/%s", CourseURL, c.Code, m.Id)
	sendRequest(http.MethodDelete, s, "", nil)
}
