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
	integrated "ash/internal/commands/managers/integrated_commands"
	"ash/internal/commands/managers/internal_actions"
	"ash/internal/configuration"
	"ash/internal/configuration/envs_loader"
	"ash/internal/dto"
	"ash/internal/executor"
	"ash/internal/internal_context"
	"ash/internal/keys_bindings"

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

	cfg := configuration.NewConfigLoader()

	// load ENVs at start
	envs_loader.LoadEnvs(cfg)

	intergratedManager := integrated.NewIntegratedManager()
	actionManager := internal_actions.NewInternalAcgionsManager()
	commandRouter := commands.NewCommandRouter(intergratedManager, actionManager)

	internalContext := internal_context.NewInternalContext(ctx, inputChan, outputChan, errs)
	promptGenerator := command_prompt.NewCommandPrompt(cfg.Prompt)
	keyBindingsManager := keys_bindings.NewKeyBindingsManager(cfg, &commandRouter)
	exec := executor.NewCommandExecutor(&commandRouter, keyBindingsManager)

	stopedChan := make(chan struct{})
	defer close(stopedChan)
	go processingInput(promptGenerator, &internalContext, exec, cfg, stopedChan)

	// waiting for stop or error xD
	e := <-errs
	internalContext.GetPrintFunction()(e.Error())
	cancelFunc()
	<-stopedChan
}

func processingInput(prompt Prompt, internalC dto.InternalContextIface, exec Executor, cfg configuration.ConfigLoader, stopedChan chan struct{}) {
	var currentBytes []byte
	ctx := internalC.GetCTX()
	outputChan := internalC.GetOutputChan()
	inputChan := internalC.GetInputChan()

	promptChan := make(chan struct{}, 1)
	defer close(promptChan)

	promptChan <- struct{}{}

mainLoop:
	for {
		select {
		case i := <-inputChan:
			if unicode.IsPrint(rune(i)) {
				currentBytes = append(currentBytes, i)
				outputChan <- i
			} else {
				if int(i) == cfg.GetKeyBind(":Backspace") {
					for _, v := range "\b\033[K" {
						outputChan <- byte(v)
					}
					currentBytes = currentBytes[:len(currentBytes)-1] // delete last input in buffer
					continue
				}
				exec.Execute(internalC.WithLastKeyPressed(i).WithCurrentInputBuffer(currentBytes))
				currentBytes = nil
				promptChan <- struct{}{}
			}
		case <-promptChan:
			prompt.GetPrompt(outputChan)
		case <-ctx.Done():
			break mainLoop
		}
	}
	stopedChan <- struct{}{}
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

func waitInterruptSignal(errs chan<- error) {
	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errs <- fmt.Errorf("%s", <-c)
	signal.Stop(c)
}
