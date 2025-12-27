package main

import (
	"log"

	"github.com/manosriram/wingman/internal/program"
)

func main() {
	targetDir := "/Users/manosriram/go/src/floppy"

	err := program.NewProgram(targetDir).Run()
	if err != nil {
		log.Fatalf("Error initializing program\n")
	}
}
