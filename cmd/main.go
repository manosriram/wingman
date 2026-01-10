package main

import (
	"log"
	"os"

	"github.com/manosriram/wingman/internal/shell"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting WorkingDir")
	}

	f, err := os.OpenFile(
		"./wingman.md",
		os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatalf("Error creating wingman.md")
	}
	defer f.Close()
	ff, err := os.OpenFile(
		"./.wingman.history.md",
		os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatalf("Error creating .wingman.history.md")
	}
	defer ff.Close()

	shell.NewShell(wd).Run()
	// shell.NewShell("/Users/manosriram/go/src/nimbusdb/").Run()
}
