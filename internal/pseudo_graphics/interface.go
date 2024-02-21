package pseudo_graphics

import "github.com/nsf/termbox-go"

type Drawer interface {
	Draw(sw PWindow, im InputManager, doneChan chan struct{}) error
}

type PWindow interface {
	Draw(x, y, w, h int, fg, bg termbox.Attribute)
	KeyInput(key rune)
	ChangeFocus()
	RemoveLastInput()
	Close()
}

type InputManager interface {
	GetInputEventChan() chan termbox.Event
}
