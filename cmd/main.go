package main

import (
	"github.com/manosriram/wingman/internal/shell"
)

func main() {
	// box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")

	// if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
	// panic(err)
	// }

	shell.NewShell().Run()

	// wd, err := os.Getwd()
	// if err != nil {
	// log.Fatalf("Error getting WorkingDir")
	// }
	// targetDir := wd

	// f, err := os.OpenFile(
	// "./wingman.md",
	// os.O_CREATE|os.O_WRONLY,
	// 0644,
	// )
	// if err != nil {
	// log.Fatalf("Error creating wingman.md")
	// }
	// defer f.Close()

	// targetDir := "/Users/manosriram/go/src/nimbusdb/"
	// // targetDir := "./"

	// err = repository.NewRepository(targetDir).Run()
	// if err != nil {
	// log.Fatalf("Error initializing program: %s\n", err.Error())
	// }
}
