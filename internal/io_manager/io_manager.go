package io_manager

import (
	"context"
	"errors"
	"io"

	"ash/internal/commands"
	"ash/internal/io_manager/list"

	"github.com/nsf/termbox-go"
)

const (
	constManagerName = "InputOutput"
	constEmptyRune   = ' '
)

type inputManager struct {
	cursorX                int
	cursorY                int
	manager                commands.CommandManagerIface
	outputCellChan         chan []termbox.Cell
	inputEventChan         chan termbox.Event
	defaultBackgroundColor termbox.Attribute
	defaultForegroundColor termbox.Attribute
}

func (i *inputManager) Read(p []byte) (n int, err error) {
	select {
	case ev := <-i.inputEventChan:
		switch ev.Type {
		case termbox.EventKey:
			if ev.Ch != 0 {
				n = copy(p, []byte{byte(ev.Ch)})
			} else {
				n = copy(p, []byte{byte(ev.Key)})
			}
			return n, nil
		}
		return 0, io.EOF
	default:
		return 0, io.EOF
	}
}

func (i *inputManager) Write(p []byte) (n int, err error) {
	for _, r := range []rune(string(p)) {
		if r == rune('\r') {
			continue
		}
		if r == rune('\t') {
			i.cursorX += 8
			continue
		}

		if r == rune('\n') {
			w, h := termbox.Size()
			i.rollScreenUp(1, w, h, termbox.GetCell, termbox.SetCell)
			i.cursorX = 0
			i.cursorY = h - 1
			continue
		}
		termbox.SetCell(i.cursorX, i.cursorY, r, i.defaultForegroundColor, i.defaultBackgroundColor)
		i.cursorX++

	}
	termbox.Flush()
	return len(p), nil
}

func (i *inputManager) Stop() {
	termbox.Interrupt()
	// termbox.Sync()
	// termbox.Flush()
}

func (i *inputManager) Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputMouse)
	return nil
}

func (i *inputManager) Start() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	defer termbox.Close()
	defer close(i.inputEventChan)
	go i.listenOutputCellChan(ctx)

	i.moveCursorAtStartPostion()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventResize:
			termbox.Sync()
			// termbox.Flush()
		case termbox.EventMouse:
			termbox.Flush()
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
	for {
		select {
		case <-ctx.Done():
			return
		case cells := <-i.outputCellChan:
			i.printCells(cells)
		}
	}
}

func (i *inputManager) printCells(cells []termbox.Cell) {
	for _, c := range cells {
		i.printSymbol(c)
	}
}

func (i *inputManager) redrawCursor() {
	termbox.SetCursor(i.cursorX, i.cursorY)
	termbox.Flush()
}

func (i *inputManager) moveCursorLeft() {
	if i.cursorX > 0 {
		i.cursorX--
	}
	i.redrawCursor()
}

func (i *inputManager) moveCursorRight() {
	i.cursorX++
	i.redrawCursor()
}

func (i *inputManager) deleteLeftSymbolAndMoveCursor() {
	if i.cursorX > 0 {
		termbox.SetCell(i.cursorX-1, i.cursorY, constEmptyRune, i.defaultForegroundColor, i.defaultBackgroundColor)
		i.moveCursorLeft()
	}
}

func (i *inputManager) moveCursorAtStartPostion() {
	i.cursorX = 0

	_, h := termbox.Size()
	i.cursorY = h - 1

	i.redrawCursor()
}

func (i *inputManager) printSymbol(c termbox.Cell) {
	if c.Ch == rune('\r') {
		return
	}
	if c.Ch == rune('\n') {
		w, h := termbox.Size()
		i.rollScreenUp(1, w, h, termbox.GetCell, termbox.SetCell)
		i.moveCursorAtStartPostion()
		return
	}
	termbox.SetCell(i.cursorX, i.cursorY, c.Ch, c.Fg, c.Bg)
	i.moveCursorRight()
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

func NewInputManager(pm promptManager) *inputManager {
	im := inputManager{
		inputEventChan:         make(chan termbox.Event),
		outputCellChan:         make(chan []termbox.Cell, 1),
		defaultBackgroundColor: termbox.ColorDefault,
		defaultForegroundColor: termbox.ColorDefault,
	}
	im.manager = commands.NewCommandManager(constManagerName, 3,
		list.NewRemoveLeftSymbol(im.deleteLeftSymbolAndMoveCursor, pm.DeleteLastSymbolFromCurrentBuffer),
	)
	return &im
}

type promptManager interface {
	DeleteLastSymbolFromCurrentBuffer() error
}

func (i *inputManager) rollScreenUp(offset, screenWidth, screenHeight int, get func(x, y int) termbox.Cell, set func(x, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute)) {
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			if y+offset >= screenHeight {
				set(x, y, constEmptyRune, i.defaultForegroundColor, i.defaultBackgroundColor)
			} else {
				curCell := get(x, y+offset)
				set(x, y, curCell.Ch, curCell.Fg, curCell.Bg)
			}
		}
	}
}