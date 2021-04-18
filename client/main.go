// This is the client program that calls the server program and displays the results
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

const (
	ContentTypeJSON = "application/json"
	CourseCode      = iota
	CourseTitle
	CourseDescription
	MaterialId
	MaterialSequence
	MaterialFileName
)

var CourseURL = "https://127.0.0.1:5000/v1/courses"

type Course struct {
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Material struct {
	Id       string `json:"id"`
	Sequence int    `json:"sequence"`
	FileName string `json:"filename"`
}

// SetupLog sets up basic logging for the program
func SetupLog() func() {
	t := time.Now()
	filename := fmt.Sprintf("log/Log%04d%02d%02d.txt", t.Year(), t.Month(), t.Day())
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("log opened")
	}
	// Besides logging, also print out to console
	wrt := io.MultiWriter(os.Stdout, file)
	log.SetOutput(wrt)
	return func() {
		err = file.Close()
		if err != nil {
			log.Fatalln(err)
		} else {
			log.Println("log closed")
		}
	}
}

// ClearConsole detects the os used and clears the console.
func ClearConsole(s string) error {
	switch s {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unable to clear console")
	}
}

// PrintMenu displays the menu to the console
func PrintMenu() {
	fmt.Println("Courses Microservice Client: ")
	fmt.Println("1. Add a new course.")
	fmt.Println("2. Add a new material to the specified course.")
	fmt.Println("3. Update a specified course.")
	fmt.Println("4. Update a specified material of a specified course.")
	fmt.Println("5. Delete a specified course.")
	fmt.Println("6. Delete a specified material.")
	fmt.Println("7. Return a list of all courses.")
	fmt.Println("8. Return a specified course.")
	fmt.Println("9. Return metadata of all materials of the specified course.")
	fmt.Println("10. Return metadata of a specified material.")
	fmt.Println("11. Download files of all materials of the specified course.")
	fmt.Println("12. Download file of a specified material.")
	fmt.Print("Select Option: ")
}

// ValidateInt converts a string to int
func ValidateInt(s *string) int {
	for {
		if _, err := fmt.Scanln(s); err != nil {
			log.Println(err)
		}
		v, err := strconv.Atoi(*s)
		if err != nil {
			log.Println(err)
		} else {
			return v
		}
	}
}

// main function starts this program
func main() {
	if os.Getenv("PLACE") == "DOCKER" {
		CourseURL = "https://server:5000/v1/courses"
	}
	logCloser := SetupLog()
	defer logCloser()
	for {
		var input string
		if err := ClearConsole(runtime.GOOS); err != nil {
			log.Println(err)
		}
		PrintMenu()
		validInput := ValidateInt(&input)
		if err := ClearConsole(runtime.GOOS); err != nil {
			log.Println(err)
		}
		switch validInput {
		case 1:
			AddNewCourseHandler()
		case 2:
			AddNewMaterialHandler()
		case 3:
			UpdateSpecificCourseHandler()
		case 4:
			UpdateSpecificMaterialHandler()
		case 5:
			RemoveSpecificCourseHandler()
		case 6:
			RemoveSpecificMaterialHandler()
		case 7:
			ReturnAllCoursesHandler()
		case 8:
			ReturnSpecificCourseHandler()
		case 9:
			ReturnAllMaterialsHandler()
		case 10:
			ReturnSpecificMaterialHandler()
		case 11:
			DownloadAllMaterialsHandler()
		case 12:
			DownloadSpecificMaterialHandler()
		default:
			log.Fatalln("invalid option entered")
		}
		if err := ClearConsole(runtime.GOOS); err != nil {
			log.Println(err)
		}
	}
}
