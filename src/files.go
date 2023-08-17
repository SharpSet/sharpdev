package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Creates default sharpdev.yml file
func genFile() {
	scriptsEx := make(map[string]string)
	scriptsEx["echo"] = "echo $_ARG1 $_ARG2"

	scriptsEx["say_secret"] = "echo SECRET"
	scriptsEx["default"] = "echo 'I am the default command'"

	values := make(map[string]string)
	values["SECRET"] = "Secret123"

	testfile := config{
		Version: 1.0,
		Scripts: scriptsEx,
		Setup:   "echo 'Setup command goes here'",
		Values:  values,
		EnvFile: ".env",
	}

	err := saveFile(testfile)
	check(err, "Failed to generate sharpdev.yml")

	fmt.Println("Created sharpdev.yml")
}

// Loads a sharpdev file
func loadFile(parent *bool) config {
	// if parent is true
	// recursively search each parent directory for a sharpdev.yml file

	var file string
	var err error
	var dir string = "./"

	if *parent {
		// find the parent directory
		dir, err = os.Getwd()
		check(err, "Failed to get current directory")
		dir = dir + "/.."

		fmt.Println("\nSearching for sharpdev.yml in parent directory")

		// loop through each directory until we find a sharpdev.yml file
		for {
			// check if sharpdev.yml file exists
			if _, err := os.Stat(dir + "/sharpdev.yml"); err == nil {
				file = dir + "/sharpdev.yml"
				break
			}

			// check if the dir leads to the root directory
			// by trying to cd to it
			err = os.Chdir(dir)
			if err != nil {
				// print err
				fmt.Println("Failed to find sharpdev.yml in parent directory")
				os.Exit(1)
			}

			// otherwise go up one directory
			dir = dir + "/.."

		}

		fmt.Println("Found sharpdev.yml in parent directory\n" + file)
		fmt.Println()
	} else {
		file = "./sharpdev.yml"
	}

	f, readErr := ioutil.ReadFile(file)

	if readErr != nil {

		file = "./env/sharpdev.yml"
		f, readErr = ioutil.ReadFile(file)

		if readErr != nil {
			fmt.Println("No sharpdev.yml was found... generating new one")
			genFile()
			return config{}
		}
	}

	var devFile config
	marshErr := yaml.Unmarshal(f, &devFile)
	if marshErr != nil {
		fmt.Println("Syntax error in sharpdev.yml")
		return config{}
	}

	if devFile.EnvFile != "" {
		devFile.EnvFile = dir + "/" + devFile.EnvFile
	}

	return devFile
}

// Saves a sharpdev file
func saveFile(devFile config) error {
	yamlData, marshErr := yaml.Marshal(devFile)
	check(marshErr, "Failed to Convert to Yaml")

	// convert yaml to string
	yamlString := string(yamlData)

	// Add newline before first occurrence of 'scripts:'
	yamlString = strings.Replace(yamlString, "scripts:", "\nscripts:", 1)
	yamlString = strings.Replace(yamlString, "values:", "\nvalues:", 1)
	yamlString = strings.Replace(yamlString, "envfile:", "\nenvfile:", 1)
	yamlString = strings.Replace(yamlString, "setup:", "\nsetup: |\n ", 1)

	// convert string to bytes
	yamlBytes := []byte(yamlString)

	writeErr := ioutil.WriteFile("./sharpdev.yml", yamlBytes, 0o644)

	if marshErr != nil || writeErr != nil {
		return errors.New("failed to save file")
	}

	return nil
}
