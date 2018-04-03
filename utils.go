package main

import (
	"io/ioutil"
	"log"
	"os"
)

func setupEnvironment(envName string, file string) error {

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Error setting environment: %v", err)
		return err
	}
	os.Setenv(envName, string(dat))
	log.Printf("%v = %v", envName, os.Getenv(envName))
	return nil
}
