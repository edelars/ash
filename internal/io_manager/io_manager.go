package io_manager

import (
	"bytes"
	"errors"
	"io"

	"ash/internal/commands"
	"ash/internal/dto"
	"ash/internal/io_manager/list"
	"ash/pkg/termbox"
)

const (
	constManagerName = "InputOutput"
	constEmptyRune   = ' '
)

type inputManager struct {
	cursorX        int
	cursorY        int
	terminateKey   uint16
	manager        commands.CommandManagerIface
	inputEventChan chan termbox.Event

	defaultBackgroundColor  termbox.Attribute
	defaultForegroundColor  termbox.Attribute
	selectedForegroundColor termbox.Attribute
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
	termbox.Sync()

	for _, r := range bytes.Runes(p) {
		if r == 13 { // /r
			continue
		}
		if r == 9 { // /t
			i.cursorX += 8
			continue
		}

		if r == 10 { // /n
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
}

func (i *inputManager) Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputMouse)
	termbox.SetOutputMode(termbox.OutputRGB)
	termbox.Clear(i.defaultForegroundColor, i.defaultBackgroundColor)
	return nil
}

func (i *inputManager) Start(execTerminateChan chan struct{}) error {
	defer termbox.Close()
	defer close(i.inputEventChan)

	i.moveCursorAtStartPostion()

mainLoop:
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
		case termbox.EventKey:
			if ev.Ch == 0 {
				switch ev.Key {
				case termbox.Key(i.terminateKey):
					execTerminateChan <- struct{}{}
					continue mainLoop
				}
			}
			i.inputEventChan <- ev
		default:
			i.inputEventChan <- ev
		}
	}
}

func (i *inputManager) redrawCursor() {
	termbox.SetCursor(i.cursorX, i.cursorY)
	termbox.SetFg(i.cursorX, i.cursorY, i.selectedForegroundColor)
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
	if c.Ch == 13 { // /r
		return
	}
	if c.Ch == 10 { // /n
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
	return func(msg string) { // TODO rewrite to Write()
		for _, c := range msg {
			i.printSymbol(termbox.Cell{Ch: c, Fg: i.defaultForegroundColor, Bg: i.defaultBackgroundColor})
		}
	}
}

func (i *inputManager) GetCellsPrintFunction() func(cells []termbox.Cell) {
	return func(cells []termbox.Cell) { // TODO rewrite to Write()
		for _, c := range cells {
			i.printSymbol(c)
		}
	}
}

func (i *inputManager) GetManager() commands.CommandManagerIface {
	return i.manager
}

func NewInputManager(pm promptManager, remSymbCmdName string, colorsAdapter dto.ColorsAdapterIface, terminateKey uint16) *inputManager {
	colors := colorsAdapter.GetColors()

	im := inputManager{
		inputEventChan:          make(chan termbox.Event),
		terminateKey:            terminateKey,
		defaultForegroundColor:  colors.DefaultForegroundColor,
		defaultBackgroundColor:  colors.DefaultBackgroundColor,
		selectedForegroundColor: colors.SelectedForegroundColor,
	}
	im.manager = commands.NewCommandManager(constManagerName, 3, false,
		list.NewRemoveLeftSymbol(remSymbCmdName, im.deleteLeftSymbolAndMoveCursor, pm.DeleteLastSymbolFromCurrentBuffer),
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
