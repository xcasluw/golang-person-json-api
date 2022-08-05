package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xcasluw/crud-go-lang/domain"
	"github.com/xcasluw/crud-go-lang/domain/person"
)

func main() {
	personService, err := person.NewService("person.json")
	if err != nil {
		fmt.Println("Error trying to create person service")
		return
	}

	http.HandleFunc("/person", func(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Not implemented", http.StatusInternalServerError)
	})

	http.ListenAndServe(":8080", nil)
}
