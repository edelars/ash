package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"ash/internal/command_prompt"
	"ash/internal/commands"
	"ash/internal/commands/managers/file_system"
	integrated "ash/internal/commands/managers/integrated_commands"
	"ash/internal/commands/managers/internal_actions"
	"ash/internal/configuration"
	"ash/internal/configuration/envs_loader"
	"ash/internal/executor"
	"ash/internal/internal_context"
	"ash/internal/io_manager"
	"ash/internal/keys_bindings"
	"ash/internal/pseudo_graphics/drawer"
)

func main() {
	errs := make(chan error, 10)
	defer close(errs)

	go waitInterruptSignal(errs)

	var wg sync.WaitGroup

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	cfg := configuration.NewConfigLoader()

	// load ENVs at start
	envs_loader.LoadEnvs(cfg)

	promptGenerator := command_prompt.NewCommandPrompt(cfg.Prompt)

	inputManager := io_manager.NewInputManager(&promptGenerator)
	if err := inputManager.Init(); err != nil {
		fmt.Println(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- inputManager.Start()
	}()

	guiDrawer := drawer.NewDrawer(cfg.GetKeyBind(":Execute"), cfg.GetKeyBind(":Close"), cfg.GetKeyBind(":Autocomplete"), cfg.GetKeyBind(":RemoveLeftSymbol"))

	// managers init
	intergratedManager := integrated.NewIntegratedManager(&cfg)
	filesystemManager := file_system.NewFileSystemManager(promptGenerator.GetUserInputFunc())
	commandRouter := commands.NewCommandRouter(intergratedManager, inputManager.GetManager(), &filesystemManager)
	actionManager := internal_actions.NewInternalActionsManager(&guiDrawer, commandRouter.GetSearchFunc(), promptGenerator.GetUserInputFunc(), cfg.Autocomplete.ShowFileInformation)
	commandRouter.AddNewCommandManager(actionManager)
	// done managers init

	internalContext := internal_context.NewInternalContext(ctx, inputManager, errs, inputManager.GetPrintFunction(), inputManager, inputManager)
	keyBindingsManager := keys_bindings.NewKeyBindingsManager(internalContext, cfg, &commandRouter)
	exec := executor.NewCommandExecutor(&commandRouter, keyBindingsManager)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- promptGenerator.Run(internalContext, &exec, cfg)
	}()

	// waiting for stop or error xD
	<-errs
	go func() {
		for range errs {
		}
	}()
	cancelFunc()

	promptGenerator.Stop()
	inputManager.Stop()

	wg.Wait()
	println("\ndone")
}

func waitInterruptSignal(errs chan<- error) {
	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errs <- fmt.Errorf("%s", <-c)
	signal.Stop(c)
}
