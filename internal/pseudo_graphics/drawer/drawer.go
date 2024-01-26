package drawer

import (
	"unicode"

	"ash/internal/pseudo_graphics"

	"github.com/nsf/termbox-go"
)

type Drawer struct {
	keyEnter       uint16
	keyClose       uint16
	keyChangeFocus uint16
	keyBackspace   uint16
	screenState    [][]termbox.Cell
}

func NewDrawer(keyEnter, keyClose, keyChangeFocus, keyBackspace uint16) Drawer {
	return Drawer{
		keyEnter:       keyEnter,
		keyClose:       keyClose,
		keyChangeFocus: keyChangeFocus,
		keyBackspace:   keyBackspace,
	}
}

func (d *Drawer) Draw(sw pseudo_graphics.PWindow, im pseudo_graphics.InputManager, doneChan chan struct{}) error {
	globalWidth, globalHeight := termbox.Size()
	termbox.Sync()
	termbox.Flush()
	d.saveScreenState(globalWidth, globalHeight, termbox.GetCell)
	defer d.restoreScreenState()

	if err := d.redrawAll(sw); err != nil {
		return err
	}
	for {

		select {
		case <-doneChan:
			return nil
		case ev := <-im.GetInputEventChan():
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.Key(d.keyChangeFocus):
					sw.ChangeFocus()
				case termbox.Key(d.keyBackspace):
					sw.RemoveLastInput()
				case termbox.Key(d.keyEnter):
					fallthrough // TODO Research data
				case termbox.Key(d.keyClose):
					return nil
				default:
					if ev.Ch != 0 && unicode.IsPrint(ev.Ch) {
						sw.KeyInput(ev.Ch) // output after select if success
					}
				}
			case termbox.EventError:
				return ev.Err
			case termbox.EventResize:
				if err := d.redrawAll(sw); err != nil {
					return err
				}
			}
		}
		if err := d.redrawAll(sw); err != nil {
			return err
		}
	}
}

func (d *Drawer) redrawAll(sw pseudo_graphics.PWindow) error {
	const coldef = termbox.ColorDefault
	if err := termbox.Sync(); err != nil {
		return err
	}
	globalWidth, globalHeight := termbox.Size()
	for x := 0; x < globalWidth; x++ {
		for y := 0; y < globalHeight; y++ {
			termbox.SetCell(x, y, ' ', coldef, coldef)
		}
	}

	if err := termbox.Clear(coldef, coldef); err != nil {
		return err
	}

	sw.Draw(0, 0, globalWidth, globalHeight)

	if err := termbox.Flush(); err != nil {
		return err
	}

	return nil
}

func (d *Drawer) saveScreenState(globalWidth, globalHeight int, get func(x, y int) termbox.Cell) {
	d.screenState = make([][]termbox.Cell, globalWidth)
	for x := 0; x < globalWidth; x++ {
		d.screenState[x] = make([]termbox.Cell, globalHeight)
		for y := 0; y < globalHeight; y++ {
			d.screenState[x][y] = get(x, y)
		}
	}
}

func (d *Drawer) restoreScreenState() {
	termbox.Sync()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	globalWidth, globalHeight := termbox.Size()
	for x := 0; x < len(d.screenState) && x < globalWidth; x++ {
		for y := 0; y < len(d.screenState[x]) && y < globalHeight; y++ {
			termbox.SetCell(x, y, d.screenState[x][y].Ch, d.screenState[x][y].Fg, d.screenState[x][y].Bg)
		}
	}
	termbox.Flush()
}
