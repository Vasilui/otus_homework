package main

import (
	"log"
	"os"
)

func main() {
	exitCode, err := ReadParams(os.Args[1:])
	if err != nil {
		log.Fatalf("failed run program: %s", err.Error())
	}

	os.Exit(exitCode)
}
