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
	flag.Parse()

	var name string

	// Load sharpdev file
	devFile, err := loadFile()
	if err != nil {
		fmt.Println("No sharpdev.yml was found... generating new one")
		genFile()
		return
	}

	if len(flag.Args()) == 0 {
		helpFunction(devFile)
		return
	}

	name = flag.Args()[0]

	err = runScript(name, devFile)
	if err != nil {
		fmt.Println(err)
	}

	return
}

func clientErrCheck(e error, msg string) {
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

func helpFunction(devFile config) {
	fmt.Println(`
This Application lets you run scripts set in your sharpdev.yml file.

It Supports:
	- env vars in the form $VAR or ${VAR}
	- Multiline commands with |
	- Inputting Args with env vars like $SHARP_ARG_{1, 2, 3, 4, etc}
	- Multiple commands can be run by using &&

Here are all the scripts you have available:
	`)

	for name := range devFile.Scripts {
		fmt.Print(name+" || ")
	}
	fmt.Println("")
}

func runScript(name string, devFile config) error {

	genSharpArgs()
	var commandStr string
	var ok bool

	if devFile.EnvFile != "" {
		err := godotenv.Load()
		clientErrCheck(err, "Failed to load env file")
	}

	if commandStr, ok = devFile.Scripts[name]; !ok {
		err := errors.New("key not in scripts config")
		clientErrCheck(err, "ScriptName "+name+" not known")
	}

	commandStrings := strings.Split(commandStr, "&&")
	for _, commStr := range commandStrings {
		err := runCommand(commStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func runCommand(commStr string) error {
	commStr, err := envsubst.String(commStr)
	clientErrCheck(err, "Failed to add ENV vars")
	arrCommandStr := strings.Fields(commStr)

	comm := arrCommandStr[0]
	args := arrCommandStr[1:]

	cmd := exec.Command(comm, args...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()

	return err
}

func genSharpArgs() {
	if len(flag.Args()) > 1 {
		for i := range flag.Args()[1:] {
			sharpArg := fmt.Sprintf("SHARP_ARG_%d", i+1)
			os.Setenv(sharpArg, flag.Args()[i+1])
		}
	}
}
