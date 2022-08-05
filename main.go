package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/xcasluw/crud-go-lang/domain"
	"github.com/xcasluw/crud-go-lang/domain/person"
)

func main() {
	personService, err := person.NewService("person.json")
	if err != nil {
		fmt.Println("Error trying to create person service")
		return
	}

	http.HandleFunc("/person/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var person domain.Person
			err := json.NewDecoder(r.Body).Decode(&person)
			if err != nil {
				fmt.Printf("Error trying to decode body. Body should be a json. Error: %s", err.Error())
				http.Error(w, "Error trying to create person", http.StatusBadRequest)
				return
			}
			if person.ID <= 0 {
				fmt.Printf("Error trying to create person. ID should be a positive integer.")
				return
			}
			// Criar pessoa
			err = personService.Create(person)
			if err != nil {
				fmt.Printf("Error trying to create person: %s", err.Error())
				http.Error(w, "Error trying to create person", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			return
		}
		if r.Method == "GET" {
			path := strings.TrimPrefix(r.URL.Path, "/person/")
			if path == "" {
				w.Header().Set("Content-type", "application/json")
				people := personService.List()
				err := json.NewEncoder(w).Encode(people)
				if err != nil {
					http.Error(w, "Error trying to list people", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			} else {
				w.Header().Set("Content-type", "application/json")
				personId, err := strconv.Atoi(path)
				if err != nil {
					http.Error(w, "Invalid id given. Person Id must be an integer", http.StatusBadRequest)
					return
				}
				person, err := personService.GetById(personId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusOK)
				err = json.NewEncoder(w).Encode(person)
				if err != nil {
					http.Error(w, "Error trying to enconde person as json", http.StatusInternalServerError)
					return
				}
			}
		}
	})

	http.ListenAndServe(":8080", nil)
}
