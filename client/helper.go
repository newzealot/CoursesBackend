package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// getInput prompts for the necessary user input to gather the data needed to perform the next step.
func getInput(items ...int) (Course, Material, error) {
	var c Course
	var m Material
	printMap := map[int]string{
		CourseCode:        "Course Code",
		CourseTitle:       "Course Title",
		CourseDescription: "Course Description",
		MaterialId:        "Material ID",
		MaterialSequence:  "Material Sequence",
		MaterialFileName:  "Material Filename",
	}

	for _, v := range items {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter %s: ", printMap[v])
		bText, _, err := reader.ReadLine()
		text := string(bText)
		if err != nil {
			log.Println(err)
		}
		switch v {
		case CourseCode:
			c.Code = text
		case CourseTitle:
			c.Title = text
		case CourseDescription:
			c.Description = text
		case MaterialId:
			m.Id = text
		case MaterialSequence:
			m.Sequence, err = strconv.Atoi(text)
			if err != nil {
				log.Println(err)
				return Course{}, Material{}, err
			}
		case MaterialFileName:
			m.FileName = text
		}

	}
	return c, m, nil
}

// sendRequest sets up and sends a http request
func sendRequest(m string, url string, ct string, b []byte) {
	// set permission for self signed
	caCert, err := ioutil.ReadFile("ssl/cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
			},
		},
	}
	// prepare the request
	req, err := http.NewRequestWithContext(context.Background(), m, url, bytes.NewBuffer(b))
	if err != nil {
		log.Println(err)
	}
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}
	if key := os.Getenv("XAPIKEY"); len(key) != 36 {
		log.Println("XAPIKEY env variable not set or not valid")
		fmt.Scanln()
		return
	} else {
		req.Header.Set("X-API-KEY", key)
	}
	// perform the request
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	log.Printf("%s: %s ... \n", m, url)
	// check if request is for downloading materials and response has a file
	if strings.Contains(url, "file") && len(res.Header["Filename"]) > 0 {
		loc := fmt.Sprintf("download/%s", res.Header["Filename"][0])
		out, err := os.Create(loc)
		if err != nil {
			log.Println(err)
		}
		defer func() {
			if err := out.Close(); err != nil {
				log.Println(err)
			}
		}()
		_, err = io.Copy(out, res.Body)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("File now in download folder")
		}
	} else {
		// receive non-file download requests
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(body))
	}
	log.Printf("STATUS: %s\n", res.Status)
	fmt.Printf("Press Enter to continue... ")
	fmt.Scanln()
}

// createMultiPart creates a MultiPart content to send both metadata and file
func createMultiPart(m Material, j []byte) (string, []byte, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	// create file metadata
	jsonWriter, err := bodyWriter.CreateFormField("json")
	if err != nil {
		log.Println("error writing to buffer")
	}
	_, err = io.Copy(jsonWriter, bytes.NewBuffer(j))
	if err != nil {
	}
	// create file itself
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", m.FileName)
	if err != nil {
		log.Println("error writing to buffer")
	}
	fh, err := os.OpenFile("upload/"+m.FileName, os.O_RDONLY, 0755)
	if err != nil {
		s := "error opening file. make sure file in upload folder"
		log.Println(s)
		fmt.Scanln()
		return "", nil, fmt.Errorf(s)
	}
	defer fh.Close()
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
	}
	// combine
	ct := bodyWriter.FormDataContentType()
	if err := bodyWriter.Close(); err != nil {
		log.Println(err)
	}
	return ct, bodyBuf.Bytes(), nil
}
