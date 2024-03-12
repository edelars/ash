package main

import (
	"ash/internal/colors_adapter"
	"ash/internal/command_prompt"
	"ash/internal/commands"
	"ash/internal/commands/managers/aliases"
	"ash/internal/commands/managers/file_system"
	"ash/internal/commands/managers/history"
	"ash/internal/commands/managers/internal_actions"
	"ash/internal/configuration"
	"ash/internal/configuration/envs_loader"
	"ash/internal/executor"
	"ash/internal/internal_context"
	"ash/internal/io_manager"
	"ash/internal/keys_bindings"
	"ash/internal/pseudo_graphics/drawer"
	"ash/internal/storage/sqlite_storage"
	"ash/internal/variables"
	"ash/pkg/escape_sequence_parser"
	"ash/version"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	integrated "ash/internal/commands/managers/integrated_commands"
)

func main() {
	errs := make(chan error, 10)
	defer close(errs)

	go waitInterruptSignal(errs)

	var wg sync.WaitGroup

	cfg := configuration.NewConfigLoader()

	// load ENVs at start
	envs_loader.LoadEnvs(cfg)

	colorsAdapter := colors_adapter.NewColorsAdapter(cfg.Colors)

	storage := sqlite_storage.NewSqliteStorage(cfg.Sqlite)
	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- storage.Run()
	}()

	promptGenerator := command_prompt.NewCommandPrompt(cfg.Prompt, colorsAdapter)

	execTerminateChan := make(chan struct{})
	defer close(execTerminateChan)

	escapeSequenceParser := escape_sequence_parser.NewEscapeSequenceParser()
	escapeSequenceDebuger := escape_sequence_parser.NewESDebug(&escapeSequenceParser, cfg.DebugOpts)

	defer escapeSequenceDebuger.Stop()

	inputManager := io_manager.NewInputManager(&promptGenerator,
		&escapeSequenceDebuger,
		configuration.CmdRemoveLeftSymbol,
		colorsAdapter,
		cfg.GetKeyBind(configuration.CmdCtrlC),
		cfg.GetKeyBind(configuration.CmdExecute),
	)
	if err := inputManager.Init(); err != nil {
		fmt.Println(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- inputManager.Start(execTerminateChan)
	}()

	guiDrawer := drawer.NewDrawer(cfg.GetKeyBind(configuration.CmdExecute),
		cfg.GetKeyBind(configuration.CmdClose),
		cfg.GetKeyBind(configuration.CmdAutocomplete),
		cfg.GetKeyBind(configuration.CmdRemoveLeftSymbol),
		colorsAdapter,
	)

	// managers init
	historyManager := history.NewHistoryManager(&storage, promptGenerator.GetUserInputFunc())
	intergratedManager := integrated.NewIntegratedManager(
		&cfg,
		version.Vesion,
		version.BuildTime,
		version.Commit,
		version.BranchName,
	)
	filesystemManager := file_system.NewFileSystemManager(promptGenerator.GetUserInputFunc())
	commandRouter := commands.NewCommandRouter(
		intergratedManager,
		inputManager.GetManager(),
		&filesystemManager,
		&historyManager,
	)
	actionManager := internal_actions.NewInternalActionsManager(&guiDrawer,
		commandRouter.GetSearchFunc(),
		promptGenerator.GetUserInputFunc(),
		cfg.Autocomplete,
		storage.SaveData,
		colorsAdapter,
	)
	commandRouter.AddNewCommandManager(actionManager)
	// done managers init

	internalContext := internal_context.NewInternalContext(inputManager, errs,
		inputManager.GetPrintFunction(),
		inputManager,
		inputManager,
		inputManager.GetCellsPrintFunction(),
		execTerminateChan,
	).WithVariables(variables.GetVariables())

	keyBindingsManager := keys_bindings.NewKeyBindingsManager(internalContext, cfg, &commandRouter)
	exec := executor.NewCommandExecutor(&commandRouter, keyBindingsManager)

	aliasManager := aliases.NewAliasesManager(cfg, &exec, cfg.GetKeyBind(configuration.CmdExecute), internalContext)
	commandRouter.AddNewCommandManager(aliasManager)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- promptGenerator.Run(internalContext, &exec, cfg, cfg.GetKeyBind(configuration.CmdExecute))
	}()

	// waiting for stop or error xD
	<-errs
	go func() {
		for range errs {
		}
	}()

	promptGenerator.Stop()
	inputManager.Stop()
	storage.Stop()
	wg.Wait()
	println("\ndone")
}

func waitInterruptSignal(errs chan<- error) {
	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errs <- fmt.Errorf("%s", <-c)
	signal.Stop(c)
}
