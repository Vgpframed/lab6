package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"log"
	"net/http"
	"os"
	"lab6/internal/structs"
	"github.com/gorilla/mux"
)

func createStudent(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var s structs.Student

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&s)

	var filename = "./" + s.Sername + "_" + s.Name + "_" + s.MiddleName + ".txt"

	wrong := analize(filename, s)
	if wrong {
		w.Write([]byte("<html> <body>Ошибка!</body></html>"))
	} else {
		writeLines(s, filename)
		ReadStrings, errRead := readLines(filename)

		if errRead != nil {
			log.Fatal(errRead)
		}
		fmt.Printf("%q\n", ReadStrings)
		for _, line := range ReadStrings {
			w.Write([]byte(line + "\n"))

		}

	}
	return
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(lines structs.Student, path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	file.WriteString(lines.Sername + "\n")
	file.WriteString(lines.Name + "\n")
	file.WriteString(lines.MiddleName + "\n")
	file.WriteString(lines.Subject + "\n")
	file.WriteString(lines.Ball + "\n")
	file.WriteString("\n")
	return file.Close()
}

func analize(filename string, s structs.Student) bool {
	var wrong bool
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	configLines := strings.Split(string(configFile), "\n")

	for i := 0; i < len(configLines)-1; i++ {

		if configLines[i] != "" {
			fmt.Println(configLines[i])
			if i > 0 {
				if configLines[i-1] == s.Subject && configLines[i] != "Неудовлетворительно" {
					wrong = true
				}
			}
			if wrong {
				break
			}
		}

	}
	fmt.Println(wrong)
	return wrong
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", createStudent)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:9000",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
