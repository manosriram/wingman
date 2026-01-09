package shell

import (
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/manosriram/wingman/internal/llm"
	"github.com/manosriram/wingman/internal/repository"
	"github.com/rivo/tview"
)

type Shell struct {
}

func NewShell() Shell {
	return Shell{}
}

func (s Shell) Run() {
	app := tview.NewApplication()

	output := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)

	input := tview.NewInputField().
		SetLabel("$ ")
	input.SetFieldBackgroundColor(tcell.ColorBlack)

	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := input.GetText()
			input.SetText("")

			if cmd == "" {
				return
			}

			fmt.Fprintf(output, "[green]$ %s\n", cmd)
			handleCommand(cmd, output)
			fmt.Fprintf(output, "\n-------------------------------------------------------------------------------------------------------------------------------------------------------\n")
		}
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(output, 0, 1, false).
		AddItem(input, 1, 0, true)

	app.SetRoot(flex, true).SetFocus(input).Run()
}

func handleCommand(line string, output *tview.TextView) {
	parts := strings.Fields(line)
	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "echo":
		fmt.Fprintf(output, "%s", args[0])
	case "help":
		fmt.Fprintf(output, "%s", "wingman")
	default:
		input := args[0]
		// call llm
		targetDir := "/Users/manosriram/go/src/nimbusdb/"
		r := repository.NewRepository(targetDir)
		err := r.Run()
		if err != nil {
			log.Fatalf("Error initializing program: %s\n", err.Error())
		}

		prompt := llm.CreateMasterPrompt(r.Signatures, input)

		// TODO: from flags
		selectedLLM := "claude"
		// selectedModel := "opus_5_2"

		response, err := llm.GetLLM(selectedLLM, prompt, llm.OPUS_4_5).Call()

		fmt.Fprintf(output, "%s", response)
	}
}
