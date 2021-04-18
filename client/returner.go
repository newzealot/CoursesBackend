package main

import (
	"fmt"
	"log"
	"net/http"
)

// ReturnAllCoursesHandler returns a list of all courses.
func ReturnAllCoursesHandler() {
	fmt.Println("7. Return a list of all courses.")
	sendRequest(http.MethodGet, CourseURL, "", nil)
}

// ReturnSpecificCourseHandler returns the specified course.
func ReturnSpecificCourseHandler() {
	fmt.Println("8. Return a specified course.")
	c, _, err := getInput(CourseCode)
	if err != nil {
		log.Println("error getting input. restarting")
		return
	}
	s := fmt.Sprintf("%s/%s", CourseURL, c.Code)
	sendRequest(http.MethodGet, s, "", nil)
}

// ReturnAllMaterialsHandler returns all materials of the specified course.
func ReturnAllMaterialsHandler() {
	fmt.Println("9. Return metadata of all materials of the specified course.")
	c, _, err := getInput(CourseCode)
	if err != nil {
		log.Println("error getting input. restarting")
		return
	}
	s := fmt.Sprintf("%s/%s/materials/metadata", CourseURL, c.Code)
	sendRequest(http.MethodGet, s, "", nil)
}

// ReturnSpecificMaterialHandler return the specified material of the specified course.
func ReturnSpecificMaterialHandler() {
	fmt.Println("10. Return metadata of a specified material.")
	c, m, err := getInput(CourseCode, MaterialId)
	if err != nil {
		log.Println("error getting input. restarting")
		return
	}
	s := fmt.Sprintf("%s/%s/materials/%s/metadata", CourseURL, c.Code, m.Id)
	sendRequest(http.MethodGet, s, "", nil)
}
