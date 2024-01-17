package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unicode"

	"ash/internal/command_prompt"
	"ash/internal/commands"
	"ash/internal/commands/managers/integrated"
	"ash/internal/configuration"
	"ash/internal/dto"
	"ash/internal/executor"
	"ash/internal/internal_context"

	"golang.org/x/term"
)

func main() {
	errs := make(chan error)
	go waitInterruptSignal(errs)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	// switch stdin into 'raw' mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	inputChan := make(chan byte)
	outputChan := make(chan byte)
	defer close(inputChan)
	defer close(outputChan)

	go readInput(ctx, inputChan)
	go writeOutput(ctx, outputChan)

	configuration.NewConfigLoader(errs)

	intergratedManager := integrated.NewIntegratedManager()
	commandRouter := commands.NewCommandRouter(intergratedManager)

	internalContext := internal_context.NewInternalContext(ctx, inputChan, outputChan, errs)
	promptGenerator := command_prompt.NewCommandPromt()
	exec := executor.NewCommandExecutor(&commandRouter)

	go processingInput(promptGenerator, &internalContext, exec)

	<-errs
	fmt.Println("ash exit")
}

func processingInput(prompt Prompt, internalC dto.InternalContextIface, exec Executor) {
	var currentBytes []byte
	ctx := internalC.GetCTX()
	outputChan := internalC.GetOutputChan()
	inputChan := internalC.GetInputChan()

	prompt.GetPrompt(outputChan)

	for {
		select {
		case i := <-inputChan:
			if unicode.IsPrint(rune(i)) {
				currentBytes = append(currentBytes, i)
				outputChan <- i
			} else {
				exec.Execute(internalC.WithLastKeyPressed(i).WithCurrentInputBuffer(currentBytes))
				currentBytes = nil
				prompt.GetPrompt(outputChan)
			}
		case <-ctx.Done():
			break
		}
	}
}

type Prompt interface {
	GetPrompt(outputChan chan byte)
}

type Executor interface {
	Execute(internalC dto.InternalContextIface)
}

func readInput(ctx context.Context, inputChan chan byte) {
	for {
		if e := ctx.Err(); e != nil {
			break
		}
		// Read the keyboad input.
		var inputArr []byte = make([]byte, 1)
		os.Stdin.Read(inputArr)
		// fmt.Printf("got new: %s\n", string(inputArr))
		inputChan <- inputArr[0]
	}
}

func writeOutput(ctx context.Context, outputChan chan byte) {
	for {
		select {
		case <-ctx.Done():
			break
		case b := <-outputChan:
			print(string(b))
		}
	}
}

// Change the directory and return the error.
// return os.Chdir(args[1])

func waitInterruptSignal(errs chan<- error) {
	fmt.Println("exit now")
	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errs <- fmt.Errorf("%s", <-c)
	signal.Stop(c)
}
