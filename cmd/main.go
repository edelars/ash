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
	"ash/internal/input_manager"
	"ash/internal/internal_context"
	"ash/internal/keys_bindings"

	"github.com/nsf/termbox-go"
)

func main() {
	errs := make(chan error)
	go waitInterruptSignal(errs)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	outputChan := make(chan byte)
	defer close(outputChan)

	go writeOutput(ctx, outputChan)

	inputManager := input_manager.NewInputManager()
	if err := inputManager.Init(); err != nil {
		fmt.Println(err)
	}

	go func() {
		errs <- inputManager.Start(ctx)
	}()

	cfg := configuration.NewConfigLoader()

	// load ENVs at start
	envs_loader.LoadEnvs(cfg)

	intergratedManager := integrated.NewIntegratedManager()
	actionManager := internal_actions.NewInternalAcgionsManager()
	commandRouter := commands.NewCommandRouter(intergratedManager, actionManager)

	internalContext := internal_context.NewInternalContext(ctx, &inputManager, outputChan, errs)
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
	inputEventChan := internalC.GetInputEventChan()

	promptChan := make(chan struct{}, 1)
	defer close(promptChan)

	promptChan <- struct{}{}

mainLoop:
	for {
		select {
		case ev := <-inputEventChan:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.Key(cfg.GetKeyBind(":Backspace")):
					for _, v := range "\b\033[K" {
						outputChan <- byte(v)
					}
					currentBytes = currentBytes[:len(currentBytes)-1] // delete last input in buffer
					continue
				default:
					if ev.Ch != 0 && unicode.IsPrint(rune(ev.Ch)) {
						currentBytes = append(currentBytes, byte(ev.Ch))
						outputChan <- byte(ev.Ch)

					} else {
						exec.Execute(internalC.WithLastKeyPressed(byte(ev.Key)).WithCurrentInputBuffer(currentBytes))
						currentBytes = nil
						promptChan <- struct{}{}
					}
				}
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
