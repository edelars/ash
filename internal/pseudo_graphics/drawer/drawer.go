package drawer

import (
	"unicode"

	"github.com/nsf/termbox-go"
)

type Drawer struct {
	keyEnter       int
	keyClose       int
	keyChangeFocus int
	keyBackspace   int
}

func NewDrawer(keyEnter, keyClose, keyChangeFocus, keyBackspace int) Drawer {
	return Drawer{
		keyEnter:       keyEnter,
		keyClose:       keyClose,
		keyChangeFocus: keyChangeFocus,
		keyBackspace:   keyBackspace,
	}
}

type inputManager interface {
	GetInputEventChan() chan termbox.Event
}

func (d Drawer) Draw(sw pWindow, im inputManager) error {
	d.redrawAll(sw)
mainloop:
	for {
		switch ev := <-im.GetInputEventChan(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.Key(d.keyChangeFocus):
				sw.ChangeFocus()
			case termbox.Key(d.keyBackspace):
				sw.RemoveLastInput()
			case termbox.Key(d.keyEnter):
				fallthrough // TODO Research data
			case termbox.Key(d.keyClose):
				break mainloop
			default:
				if ev.Ch != 0 && unicode.IsPrint(ev.Ch) {
					sw.KeyInput(ev.Ch) // output after select if success
				}
			}
		case termbox.EventError:
			return ev.Err
		case termbox.EventResize:
			d.redrawAll(sw)
		}

		d.redrawAll(sw)
	}

	return nil
}

type pWindow interface {
	Draw(x, y, w, h int)
	KeyInput(key rune)
	ChangeFocus()
	RemoveLastInput()
}

func (d Drawer) redrawAll(sw pWindow) {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	globalWidth, globalHeight := termbox.Size()
	sw.Draw(0, 0, globalWidth, globalHeight)

	termbox.Flush()
}
