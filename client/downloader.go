package main

import (
	"fmt"
	"log"
	"net/http"
)

// DownloadAllMaterialsHandler downloads all materials in the form of a zip file.
// It is stored in the download folder
func DownloadAllMaterialsHandler() {
	fmt.Println("11. Download files of all materials of the specified course.")
	c, _, err := getInput(CourseCode)
	if err != nil {
		log.Println("error getting input. restarting")
		return
	}
	s := fmt.Sprintf("%s/%s/materials/files", CourseURL, c.Code)
	sendRequest(http.MethodGet, s, "", nil)
}

// DownloadSpecificMaterialHandler downloads the specified material.
// It is stored in the download folder.
func DownloadSpecificMaterialHandler() {
	fmt.Println("12. Download file of a specified material.")
	c, m, err := getInput(CourseCode, MaterialId)
	if err != nil {
		log.Println("error getting input. restarting")
		return
	}
	s := fmt.Sprintf("%s/%s/materials/%s/files", CourseURL, c.Code, m.Id)
	sendRequest(http.MethodGet, s, "", nil)
}
