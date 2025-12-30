package main

import (
	"log"

	"github.com/manosriram/wingman/internal/repository"
)

func main() {
	// wd, err := os.Getwd()
	// if err != nil {
	// log.Fatalf("Error getting WorkingDir")
	// }
	// targetDir := wd

	// targetDir := "/Users/manosriram/go/src/syncthing/"
	targetDir := "./"

	err := repository.NewRepository(targetDir).Run()
	if err != nil {
		log.Fatalf("Error initializing program: %s\n", err.Error())
	}
}
