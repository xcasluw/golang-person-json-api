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
				// List all
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				people := personService.List()
				err := json.NewEncoder(w).Encode(people)
				if err != nil {
					http.Error(w, "Error trying to list people", http.StatusInternalServerError)
					return
				}
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
		if r.Method == "PUT" {
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
			// Atualizar pessoa
			err = personService.Update(person)
			if err != nil {
				fmt.Printf("Error trying to update person: %s", err.Error())
				http.Error(w, "Error trying to update person", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == "DELETE" {
			path := strings.TrimPrefix(r.URL.Path, "/person/")
			if path == "" {
				http.Error(w, "ID must be provided", http.StatusBadRequest)
				return
			}

			personId, err := strconv.Atoi(path)
			if err != nil {
				http.Error(w, "Invalid id given. Person Id must be an integer", http.StatusBadRequest)
				return
			}
			err = personService.DeleteById(personId)
			if err != nil {
				fmt.Printf("Error trying to delete person: %s", err.Error())
				http.Error(w, "Error trying to delete person", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
