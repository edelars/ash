package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	errs := make(chan error)
	go waitInterruptSignal(errs)

	// intergratedManager := integrated.NewIntegratedManager()
	// commandRouter := commands.NewCommandRouter(intergratedManager)

	inChan := make(chan []byte)
	go readInput(inChan)

	for {
		select {
		case i := <-inChan:
			fmt.Println(string(i))
		case <-errs:
			break
		}
	}
}

func readInput(outChan chan []byte) {
	reader := bufio.NewReader(os.Stdin)
	for {
		// fmt.Print("> ")
		// Read the keyboad input.
		input, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Printf("got: %s\n", string(input))
		outChan <- input
	}
}

// ErrNoPath is returned when 'cd' was called without a second argument.
var ErrNoPath = errors.New("path required")

func execInput(input string) error {
	// Split the input separate the command and the arguments.
	args := strings.Split(input, " ")

	// Check for built-in commands.
	switch args[0] {
	case "cd":
		// 'cd' to home with empty path not yet supported.
		if len(args) < 2 {
			return ErrNoPath
		}
		// Change the directory and return the error.
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	}

	// Prepare the command to execute.
	cmd := exec.Command(args[0], args[1:]...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command return the error.
	return cmd.Run()
}

func waitInterruptSignal(errs chan<- error) {
	fmt.Println("exit now")
	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errs <- fmt.Errorf("%s", <-c)
	signal.Stop(c)
}
