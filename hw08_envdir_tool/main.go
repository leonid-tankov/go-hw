package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Printf("Error: not enough args\n")
		return
	}
	environment, err := ReadDir(args[1])
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}
	code := RunCmd(args[2:], environment)
	os.Exit(code)
}
