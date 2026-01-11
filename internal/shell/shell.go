package shell

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/manosriram/wingman/internal/llm"
	"github.com/manosriram/wingman/internal/repository"
	"github.com/rivo/tview"
)

type ProgramFlags struct {
	Model *string
}

type Shell struct {
	ShellDir   string
	Flags      ProgramFlags
	Repository *repository.Repository
	LLM        llm.LLM
}

func NewShell(targetDir string) (Shell, error) {
	modelPtr := flag.String("model", "", "Model of the LLM")
	flag.Parse()
	llm, err := llm.NewLLM(*modelPtr)
	if err != nil {
		return Shell{}, err
	}

	return Shell{
		Flags: ProgramFlags{
			Model: modelPtr,
		},
		ShellDir: targetDir,
		LLM:      llm,
	}, nil
}

type CmdChannel struct {
	Response string
	Error    error
}

func (s Shell) Run() {
	app := tview.NewApplication()

	// targetDir := "/Users/manosriram/go/src/nimbusdb/"
	r := repository.NewRepository(s.ShellDir)
	err := r.Run()
	if err != nil {
		log.Fatalf("Error initializing program: %s\n", err.Error())
	}
	s.Repository = r

	output := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)

	input := tview.NewInputField().
		SetLabel("$ ")
	input.SetFieldBackgroundColor(tcell.ColorBlack)

	ch := make(chan (CmdChannel))

	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := input.GetText()
			input.SetText("")

			if cmd == "" {
				return
			}

			fmt.Fprintf(output, "[green]$ %s\n", cmd)

			// fmt.Fprintf(output, "%s\n", cmd)

			go func() {
				go s.handleCommand(cmd, ch, output)
				result := <-ch
				if result.Error != nil {
					// log.Fatalf("Error %s", result.Error.Error())
					fmt.Fprintf(output, "%s\n", result.Error.Error())
				}
				app.QueueUpdateDraw(func() {
					fmt.Fprintf(output, "%s\n", result.Response)
					fmt.Fprintf(output, "\n-------------------------------------------------------------------------------------------------------------------------------------------------------\n")
					output.ScrollToEnd()
				})
			}()
		}
	})

	// Allow Tab key to switch focus to output for scrolling
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(output)
			return nil
		}
		return event
	})

	// Allow Tab/Enter to return focus to input
	output.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyEnter {
			app.SetFocus(input)
			return nil
		}
		return event
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(output, 0, 1, false).
		AddItem(input, 1, 0, true)

	app.SetRoot(flex, true).SetFocus(input).EnableMouse(true).Run()
}

func (s Shell) handleCommand(line string, ch chan<- CmdChannel, output *tview.TextView) {
	parts := strings.Fields(line)
	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "echo":
		fmt.Fprintf(output, "%s", args[0])
	case "help":
		fmt.Fprintf(output, "%s", "wingman")
	case "/clear":
		output.Clear()
	case "/exit":
		os.Exit(0)
	case "/add":
		paths := args
		err := s.Repository.AddFiles(paths)
		if err != nil {
			fmt.Fprintf(output, "%s", "Error adding file(s): "+err.Error()+"\n")
		} else {
			fmt.Fprintf(output, "%s", "Added file(s)\n")
		}
	case "/drop":
		paths := args
		s.Repository.DropFiles(paths)
		fmt.Fprintf(output, "%s", "Dropped file(s)\n")
	default:
		cmdCh := CmdChannel{}

		input := strings.Join(parts, " ")
		prompt := s.Repository.CreateMasterPrompt(input)
		response, err := s.LLM.Call(prompt)
		if err != nil {
			cmdCh.Error = err
		} else {
			cmdCh.Response = response.Response
		}

		ch <- cmdCh

		err = s.LLM.WriteToHistory(input, response)
		if err != nil {
			return // TODO: Handle err
		}

	}
}
