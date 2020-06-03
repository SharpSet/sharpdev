package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/a8m/envsubst"
	"github.com/joho/godotenv"
)

func main() {
	var name string

	// Load sharpdev file
	devFile, err := loadFile()
	if err != nil {
		fmt.Println("No sharpdev.yml was found... generating new one")
		genFile()
		return
	}

	// Make Helper Function and Parse Flags
	setHelperFunction(devFile)
	flag.Parse()

	// If no script is called load helpfunction
	if len(flag.Args()) == 0 {
		flag.Usage()
		return
	}
	// Run script with name of first arg
	name = flag.Args()[0]
	err = runScript(name, devFile)
	if err != nil {
		fmt.Println(err)
	}

	return
}

// Deals with client Errors
func check(e error, msg string) {
	// Try and get SHARPDEV var
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cannot read enviroment")
	}

	if e != nil {
		if os.Getenv("SHARPDEV") == "TRUE" {
			fmt.Println(e)
		}
		log.Fatal(msg)
	}
}

func setHelperFunction(devFile config) {
	flag.Usage = func() {
		fmt.Println(`
This Application lets you run scripts set in your sharpdev.yml file.

It Supports:
	- env vars in the form $VAR or ${VAR}
	- Multiline commands with |
	- Inputting Args with env vars like $SHARP_ARG_{1, 2, 3, 4, etc}

Here are all the scripts you have available:
			`)

		// Shows all script name
		for name := range devFile.Scripts {
			fmt.Print(name + " || ")
		}
		fmt.Println("")
	}
}

func runScript(name string, devFile config) error {

	// Create Env Vars from other args
	genSharpArgs()
	var commandStr string
	var ok bool

	// Check if a envfile is required
	if devFile.EnvFile != "" {
		err := godotenv.Load()
		check(err, "Failed to load env file")
	}

	// Check that the arg is actually a script
	if commandStr, ok = devFile.Scripts[name]; !ok {
		err := errors.New("key not in scripts config")
		check(err, "ScriptName "+name+" not known")
	}

	// For each command in a script split by &&
	commandStrings := strings.Split(commandStr, "&&")
	for _, commStr := range commandStrings {

		// Run command
		err := runCommand(commStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func runCommand(commStr string) error {
	// Substiute Env Vars
	commStr, err := envsubst.String(commStr)
	check(err, "Failed to add ENV vars")
	arrCommandStr := strings.Fields(commStr)

	comm := arrCommandStr[0]
	args := arrCommandStr[1:]

	// Run command through OS args
	cmd := exec.Command(comm, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()

	return err
}

func genSharpArgs() {

	// If there is more than one arg
	if len(flag.Args()) > 1 {
		for i := range flag.Args()[1:] {

			// Add arg to Environ
			sharpArg := fmt.Sprintf("SHARP_ARG_%d", i+1)
			os.Setenv(sharpArg, flag.Args()[i+1])
		}
	}
}
