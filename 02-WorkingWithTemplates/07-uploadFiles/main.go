package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("error getting a file for the provided form key: %v", err)
		return
	}
	defer file.Close()

	//*We are going to use path.Join() function for avoiding the error when finding tmp directory
	//! You must create your own tmp dir in order to make this code works
	wd, _ := os.Getwd() //This line gets the working directory path
	out, pathError := os.Create(path.Join(wd, "/tmp/uploadedFile"))
	if pathError != nil {
		log.Printf("error creating a file for writing: %v", pathError)
		return
	}
	defer out.Close()
	_, copyFileError := io.Copy(out, file)
	if copyFileError != nil {
		log.Printf("error ocurred while file copy: %v", copyFileError)
		return
	}
	fmt.Fprintf(w, "File uploaded successfully: "+header.Filename)
}

func index(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("templates/upload-file.html")
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		log.Printf("Error while executing the template: %v", err)
		return
	}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", fileHandler)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
		return
	}
}
