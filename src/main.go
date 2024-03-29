package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

// add a -p flag var as bool
var (
	parent    = flag.Bool("p", false, "Use a parent sharpdev.yml file")
	version   = flag.Bool("v", false, "Get the version number")
	version2  = flag.Bool("version", false, "Get the version number")
	skipSetup = flag.Bool("ss", false, "Skips using the setup option")
	dotFile   = flag.String("url", "", "The dotfile repo for the sharpdev files. See https://github.com/Sharpz7/dotfiles")
	envName   = flag.String("envname", "", "The name of sharpdev env to download from")
)

func main() {
	var name string

	flag.Parse()

	// If -v is used, print version
	if *version || *version2 {
		fmt.Println(Version)
		os.Exit(0)
	}

	// if dotfile is set
	if *dotFile != "" {
		// Place dotFile/env/envName/sharpdev.yml in ./env
		err := downloadDotFile(*dotFile, *envName)

		check(err, "Failed to download dotfile")
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
}

// Deals with client Errors
func check(e error, msg string) {
	// Try and get SHARPDEV var
	godotenv.Load()

	if e != nil {
		if os.Getenv("DEV") == "TRUE" {
			fmt.Println(e)
		}
		log.Fatal(msg + "\n" + e.Error())
	}
}

func setHelperFunction(devFile config) {
	flag.Usage = func() {
		fmt.Println(`
This Application lets you run scripts set in your sharpdev.yml file.

Note that if no file is found in the dir you are in, it will instead search in ./env

It Supports:
	- env vars in the form $VAR or ${VAR}
	- Multiline commands with |
	- Inputting Args with env vars like $_ARG{1, 2, 3, 4, etc}

Flags:
	-p Uses a parent sharpdev.yml file

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

	// Run command
	err = runCommand(commandStr, devFile)
	if err != nil {
		return err
	}

	return nil
}

func runCommand(commStr string, devFile config) error {
	if devFile.Setup != "" && !*skipSetup {
		// add setup command to commStr
		commStr = devFile.Setup + "\n" + commStr
	}

	// Get Input Args
	commStr = placeInputArgs(commStr)

	// For command string replace any reference to args
	for key, val := range devFile.Values {
		commStr = strings.ReplaceAll(commStr, key, val)
	}

	// Run command through OS args
	cmd := exec.Command("/bin/sh", "-c", commStr)

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	return err
}

func placeInputArgs(commStr string) string {
	if len(flag.Args()) == 0 {
		return commStr
	}

	var previousCommStr string

	// If there is more than one arg
	if len(flag.Args()) > 1 {
		for i := range flag.Args()[1:] {

			// Add arg to CommStr
			commStr = strings.ReplaceAll(commStr, fmt.Sprintf("$_ARG%d", i+1), flag.Args()[i+1])

			// check if commStr changed from previous iteration
			if commStr == previousCommStr {
				// convert all args to string with spaces
				args := strings.Join(flag.Args()[i+1:], " ")
				fmt.Println("Extra args were provided - appending to the end (" + args + ")")
				return commStr + " " + args
			}

			previousCommStr = commStr
		}
	}

	return commStr
}

func checkVersion(devFile config) error {
	if devFile.Version != 1.0 {
		return errors.New("")
	}

	return nil
}
