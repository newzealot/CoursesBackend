package main

// PUT /v1/courses/{code} – Update the specified course.
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// UpdateSpecificCourseHandler updates the specified course.
// The COURSE_ID cannot be changed. If there is a need to do so, delete and add course again.
func UpdateSpecificCourseHandler() {
	fmt.Println("3. Update a specified course.")
	c, _, err := getInput(CourseCode, CourseTitle, CourseDescription)
	if err != nil {
		log.Println("error getting input. restarting")
		return
	}
	s := fmt.Sprintf("%s/%s", CourseURL, c.Code)
	requestBody, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Println(err)
	}
	sendRequest(http.MethodPut, s, ContentTypeJSON, requestBody)
}

// UpdateSpecificMaterialHandler updates the specified material in the specified course.
// The MATERIAL_ID cannot be changed. If there is a need to do so, delete and add material again.
func UpdateSpecificMaterialHandler() {
	fmt.Println("Update a specified material of a specified course.")
	c, m, err := getInput(CourseCode, MaterialId, MaterialSequence, MaterialFileName)
	if err != nil {
		log.Println("error getting input. restarting")
		return
	}
	s := fmt.Sprintf("%s/%s/materials/%s", CourseURL, c.Code, m.Id)
	j, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		log.Println(err)
	}
	ct, b, err := createMultiPart(m, j)
	if err != nil {
		return
	}
	sendRequest(http.MethodPut, s, ct, b)
}
