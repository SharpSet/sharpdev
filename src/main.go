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

// add a -p flag var as bool
var parent = flag.Bool("p", false, "Use a parent sharpdev.yml file")
var version = flag.Bool("v", false, "Get the version number")
var version2 = flag.Bool("version", false, "Get the version number")

func main() {
	var name string

	flag.Parse()

	// If -v is used, print version
	if *version || *version2 {
		fmt.Println(Version)
		os.Exit(0)
	}

	// Load sharpdev file
	devFile := loadFile(parent)
	if devFile.Version == 0 {
		os.Exit(1)
	}

	// Check if a envfile is required
	if devFile.EnvFile != "" {
		err := godotenv.Load(devFile.EnvFile)
		check(err, "Failed to load env file2")
	}

	// Make Helper Function and Parse Flags
	setHelperFunction(devFile)

	if len(flag.Args()) == 0 {
		name = "default"
	} else if flag.Args()[0] == "help" {
		flag.Usage()
		return
	} else {
		name = flag.Args()[0]
	}

	// Run script with name of first arg
	err := runScript(name, devFile)
	if err != nil {
		fmt.Println(err)
	}

	return
}

// Deals with client Errors
func check(e error, msg string) {
	// Try and get SHARPDEV var
	godotenv.Load()

	if e != nil {
		if os.Getenv("DEV") == "TRUE" {
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
	- Inputting Args with env vars like $_ARG{1, 2, 3, 4, etc}

Flags:
	-p  Uses a parent sharpdev.yml file

If no script is called, the "default" script will be run.

Here are all the scripts you have available:
			`)

		// Shows all script name
		for name := range devFile.Scripts {
			if name != "default" {
				fmt.Print(name + " || ")
			}
		}
		fmt.Println("")
	}
}

func runScript(name string, devFile config) error {

	// Check if version is correct
	err := checkVersion(devFile)
	check(err, "Incorrect version. \nCurrently running 1.0, Script is running "+fmt.Sprint(devFile.Version))

	// Create Env Vars from other args
	genSharpArgs()
	var commandStr string
	var ok bool

	// Check that the arg is actually a script
	if commandStr, ok = devFile.Scripts[name]; !ok {
		err := errors.New("key not in scripts config")

		if name == "default" {
			check(err, "There is no default script set")
		} else {
			check(err, "ScriptName "+name+" not known")
		}

	}

	// Run Setup
	setupCommand := devFile.Setup
	err = runCommand(setupCommand, devFile)
	if err != nil {
		return err
	}

	// Run command
	err = runCommand(commandStr, devFile)
	if err != nil {
		return err
	}

	return nil
}

func runCommand(commStr string, devFile config) error {
	// Substitute Env Vars
	commStr, err := envsubst.String(commStr)
	check(err, "Failed to add ENV vars")

	// For command string replace any reference to args
	for key, val := range devFile.Values {
		commStr = strings.ReplaceAll(commStr, key, val)
	}

	// Replace "\n" with &&
	strings.Replace(commStr, "\n", "&&", -1)

	// Run command through OS args
	cmd := exec.Command("/bin/sh", "-c", commStr)
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
			sharpArg := fmt.Sprintf("_ARG%d", i+1)
			os.Setenv(sharpArg, flag.Args()[i+1])
		}
	}
}

func checkVersion(devFile config) error {

	if devFile.Version != 1.0 {
		return errors.New("")
	}

	return nil
}
