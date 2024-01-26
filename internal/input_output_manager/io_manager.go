package input_manager

import (
	"context"
	"errors"

	"ash/internal/commands"
	"ash/internal/input_output_manager/list"

	"github.com/nsf/termbox-go"
)

const constManagerName = "InputOutput"

type inputManager struct {
	cursorX                int
	cursorY                int
	manager                commands.CommandManagerIface
	outputCellChan         chan []termbox.Cell
	inputEventChan         chan termbox.Event
	defaultBackgroundColor termbox.Attribute
	defaultForegroundColor termbox.Attribute
	currentUserInput       []termbox.Cell
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

func (i *inputManager) listenOutputCellChan(ctx context.Context) {
	defer close(i.outputCellChan)
	select {
	case <-ctx.Done():
		return
	case cells := <-i.outputCellChan:
		i.printCells(cells)
	}
}

func (i *inputManager) printCells(cells []termbox.Cell) {
	for _, c := range cells {
		i.printSymbol(c)
	}
}

func (i *inputManager) redrawCursor() {
	termbox.SetCursor(i.cursorX, i.cursorY)
}

func (i *inputManager) moveCursorLeft() {
	i.cursorX--
	i.redrawCursor()
}

func (i *inputManager) moveCursorRigth() {
	i.cursorX++
	i.redrawCursor()
}

func (i *inputManager) deleteLeftSymbolAndMoveCursor() {
	i.moveCursorLeft()
}

func (i *inputManager) moveCursorAtStartPostion() {
	i.moveCursorLeft()
}

func (i *inputManager) printSymbol(c termbox.Cell) {
	if c.Ch == rune('\n') {
		rollScreenUp(termbox.GetCell, termbox.SetCell)
		i.currentUserInput = nil
		i.moveCursorAtStartPostion()
		return
	}

	termbox.SetCell(i.cursorX, i.cursorY, c.Ch, c.Fg, c.Bg)
	i.moveCursorRigth()
}

func (i *inputManager) GetInputEventChan() chan termbox.Event {
	return i.inputEventChan
}

func (i *inputManager) GetPrintFunction() func(msg string) {
	return func(msg string) {
		var r []termbox.Cell
		for _, c := range msg {
			r = append(r, termbox.Cell{Ch: c, Fg: i.defaultForegroundColor, Bg: i.defaultBackgroundColor})
		}
		i.outputCellChan <- r
	}
}

func (i *inputManager) GetManager() commands.CommandManagerIface {
	return i.manager
}

func NewInputManager() inputManager {
	im := inputManager{
		inputEventChan:         make(chan termbox.Event),
		outputCellChan:         make(chan []termbox.Cell),
		defaultBackgroundColor: termbox.ColorDefault,
		defaultForegroundColor: termbox.ColorDefault,
	}
	im.manager = commands.NewCommandManager(constManagerName,
		list.NewRemoveLeftSymbol(im.deleteLeftSymbolAndMoveCursor),
	)

	return im
}

func rollScreenUp(get func(x, y int) termbox.Cell, set func(x, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute)) {
}
