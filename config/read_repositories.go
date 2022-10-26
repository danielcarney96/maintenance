package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Requirement struct {
	Key     string
	Version string
}

type Repo struct {
	Url          string
	Requirements []Requirement
}

func ReadRepositoriesFromFile(filename string) map[string]Repo {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	m := make(map[string]Repo)

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatal(err)
	}

	return m
}
