package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func setupEnvironment(pathScript string) error {
	f, err := ioutil.ReadFile(pathScript)
	if err != nil {
		log.Printf("Error opening file %v. Error: %v ", pathScript, err)
		return err
	}

	ex := `(?m)^export (.*)="(.*)"$`
	r := regexp.MustCompile(ex)
	vars := r.FindAllStringSubmatch(string(f), -1)
	if vars == nil {
		message := fmt.Sprintf("No env vars found with pattern %v", ex)
		log.Printf(message)
		return errors.New(message)
	}

	for _, v := range vars {
		value := strings.Replace(v[2], `\`, "", -1) // Remove escape chars
		log.Printf("Defining %v=%v", v[1], value)
		os.Setenv(v[1], value)
	}

	return nil
}
