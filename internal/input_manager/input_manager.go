package input_manager

import (
	"context"
	"errors"

	"github.com/nsf/termbox-go"
)

type inputManager struct {
	inputChan      chan byte
	inputEventChan chan termbox.Event
}

func (i *inputManager) Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputEsc)
	return nil
}

func (i *inputManager) Start(ctx context.Context) error {
	defer termbox.Close()
	defer close(i.inputEventChan)

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		switch ev := termbox.PollEvent(); ev.Type {

		case termbox.EventError:
			return ev.Err
		case termbox.EventInterrupt:
			return errors.New("got EventInterrupt, exiting")
		default:
			i.inputEventChan <- ev
		}
	}
}

func (i *inputManager) GetInputEventChan() chan termbox.Event {
	return i.inputEventChan
}

func NewInputManager() inputManager {
	return inputManager{
		inputEventChan: make(chan termbox.Event),
	}
}
