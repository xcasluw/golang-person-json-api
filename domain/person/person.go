package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xcasluw/crud-go-lang/domain"
)

type Service struct {
	dbFilePath string
	people     domain.People
}

func NewService(dbFilePath string) (Service, error) {
	// Verifica se o arquivo existe
	_, err := os.Stat(dbFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Caso não exista, cria um arquivo vazio
			err = createEmptyFile(dbFilePath)
			if err != nil {
				return Service{}, err
			}
			return Service{
				dbFilePath: dbFilePath,
				people:     domain.People{},
			}, nil
		}
	}

	// Se existir, lê o arquivo e atualiza a variável people do serviço com as pessoas do arquivo
	jsonFile, err := os.Open(dbFilePath)
	if err != nil {
		return Service{}, fmt.Errorf("Error trying to open file that contains all people: %s", err.Error())
	}

	jsonFileContentByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Service{}, fmt.Errorf("Error trying to read file: %s", err.Error())
	}

	var allPeople domain.People
	json.Unmarshal(jsonFileContentByte, &allPeople)

	return Service{
		dbFilePath: dbFilePath,
		people:     allPeople,
	}, nil
}

func createEmptyFile(dbFilePath string) error {
	var people domain.People = domain.People{
		People: []domain.Person{},
	}
	peopleJSON, err := json.Marshal(people)
	if err != nil {
		return fmt.Errorf("Error trying to encode people as JSON?: %s", err.Error())
	}

	err = ioutil.WriteFile(dbFilePath, peopleJSON, 0755)
	if err != nil {
		return fmt.Errorf("Error trying to write to file. Error: %s", err.Error())
	}

	return nil
}

func (s *Service) Create(person domain.Person) error {
	// Verifica se a pessoa já existe
	if s.exists(person) {
		return fmt.Errorf("Error trying to create person. Peson with ID exists")
	}

	// Adiciona a pessoa na slice de pessoas
	s.people.People = append(s.people.People, person)

	// Salva o arquivo
	err := s.saveFile()
	if err != nil {
		return fmt.Errorf("Error trying saving file in method Create. Error: %s", err.Error())
	}

	return nil
}

func (s Service) exists(person domain.Person) bool {
	for _, currentPerson := range s.people.People {
		if currentPerson.ID == person.ID {
			return true
		}
	}
	return false
}

func (s Service) saveFile() error {
	allPeopleJSON, err := json.Marshal(s.people)
	if err != nil {
		return fmt.Errorf("Error trying to encode people as json %s", err.Error())
	}
	return ioutil.WriteFile(s.dbFilePath, allPeopleJSON, 0755)
}

func (s Service) List() domain.People {
	return s.people
}

func (s Service) GetById(personID int) (domain.Person, error) {
	for _, currentPerson := range s.people.People {
		if currentPerson.ID == personID {
			return currentPerson, nil
		}
	}
	return domain.Person{}, fmt.Errorf("Person not found")
}
