package selection_window

import (
	"ash/internal/dto"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	constHelpMessage   = "Press <tab> to change focus and <[a-z]> to execute command from search results"
	constNoDataMessage = "No search results"
)

type selectionWindow struct {
	cursorX int
	cursorY int
	mainX   int
	mainY   int
	mainW   int
	mainH   int

	defaultBackgroundColor termbox.Attribute
	defaultForegroundColor termbox.Attribute
	sourceBackgroundColor  termbox.Attribute
	sourceForegroundColor  termbox.Attribute
	srKeyBackgroundColor   termbox.Attribute
	srKeyForegroundColor   termbox.Attribute

	focused      bool
	dataSources  dto.DataSource
	symbolsMap   map[rune]rune
	currentInput []rune
	searchFunc   func(patter []rune) dto.DataSource
	resultFunc   func(cmd dto.CommandIface, userInput []rune)
}

func NewSelectionWindow(userInput []rune, searchFunc func(patter []rune) dto.DataSource, resultFunc func(cmd dto.CommandIface, userInput []rune)) selectionWindow {
	sw := selectionWindow{
		defaultBackgroundColor: 0,
		defaultForegroundColor: 0,
		sourceBackgroundColor:  termbox.ColorLightBlue,
		sourceForegroundColor:  termbox.ColorBlack,
		srKeyBackgroundColor:   termbox.ColorLightGreen,
		srKeyForegroundColor:   termbox.ColorBlack,
		focused:                false,
		symbolsMap:             map[rune]rune{},
		currentInput:           userInput,
		searchFunc:             searchFunc,
		resultFunc:             resultFunc,
	}
	if runewidth.EastAsianWidth {
		sw.symbolsMap = map[rune]rune{'─': '-', '│': '|', '┌': '+', '└': '+', '┐': '+', '┘': '+'}
	} else {
		sw.symbolsMap = map[rune]rune{'─': '─', '│': '│', '┌': '┌', '└': '└', '┐': '┐', '┘': '┘'}
	}
	return sw
}

func (sw *selectionWindow) updateDataSource() {
	sw.dataSources = sw.searchFunc(sw.currentInput)
}

func (sw *selectionWindow) KeyInput(key rune) {
	if sw.focused {
		sw.currentInput = append(sw.currentInput, rune(key))
		sw.reDraw(true)
	} else {
		sw.commandChoosed(key)
	}
}

func (sw *selectionWindow) commandChoosed(key rune) {
	if cmd := sw.dataSources.GetCommand(key); cmd != nil {
		sw.resultFunc(cmd, sw.currentInput)
	}
}

func (sw *selectionWindow) ChangeFocus() {
	sw.focused = !sw.focused
	sw.drawCursor()
}

func (sw *selectionWindow) setCursor(x, y int) {
	sw.cursorX = x
	sw.cursorY = y
	sw.drawCursor()
}

func (sw *selectionWindow) drawCursor() {
	if sw.focused {
		termbox.SetCursor(sw.cursorX, sw.cursorY)
	} else {
		termbox.HideCursor()
	}
}

func (sw *selectionWindow) reDraw(clearScreen bool) {
	if clearScreen {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	}
	const bottomInputBarH = 2
	const overheadSpaceForSource = 2

	sw.updateDataSource()

	dataResult := sw.dataSources.GetData(sw.mainH-bottomInputBarH, overheadSpaceForSource)

	switch len(dataResult) {
	case 0: // if no data we will draw empty results message
		x := sw.mainW/2 - len([]rune(constNoDataMessage))/2 // calculate center of screen
		y := (sw.mainH - bottomInputBarH) / 2
		tbprint(x, y, sw.defaultForegroundColor, sw.defaultBackgroundColor, constNoDataMessage)
	default: // draw data
		stepY := sw.mainY
		for _, ds := range dataResult {
			stepY = sw.drawSource(sw.mainX+1, stepY, sw.mainW, ds)
		}
	}

	// draw help
	tbprint(sw.mainX+1, sw.mainY+sw.mainH-4, sw.defaultForegroundColor, sw.defaultBackgroundColor, constHelpMessage)
	// draw input
	sw.drawRectangle(sw.mainX, sw.mainY+sw.mainH-3, sw.mainW, 3) // draw input box
	curX := sw.mainX + 1
	curY := sw.mainY + sw.mainH - 2

	tbprint(curX, curY, termbox.ColorDefault, termbox.ColorDefault, string(sw.currentInput))
	curX = curX + len(sw.currentInput)
	sw.setCursor(curX, curY)
}

func (sw *selectionWindow) Draw(x, y, w, h int) {
	sw.mainX = x
	sw.mainY = y
	sw.mainW = w
	sw.mainH = h
	sw.reDraw(false)
}

func (sw *selectionWindow) fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, sw.defaultForegroundColor, sw.defaultBackgroundColor)
		}
	}
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += 1
	}
}

func (sw *selectionWindow) drawRectangle(x, y, w, h int) {
	const coldef = termbox.ColorDefault

	sw.fill(x+1, y, w-2, 1, termbox.Cell{Ch: sw.getDrawSymbol('─')})     // top line
	sw.fill(x+1, y+h-1, w-2, 1, termbox.Cell{Ch: sw.getDrawSymbol('─')}) // bottom line

	sw.fill(x, y+1, 1, h-2, termbox.Cell{Ch: sw.getDrawSymbol('│')})     // left line
	sw.fill(x+w-1, y+1, 1, h-2, termbox.Cell{Ch: sw.getDrawSymbol('│')}) // right line

	termbox.SetCell(x, y, sw.getDrawSymbol('┌'), sw.defaultForegroundColor, sw.defaultForegroundColor)
	termbox.SetCell(x, y+h-1, sw.getDrawSymbol('└'), sw.defaultForegroundColor, sw.defaultForegroundColor)
	termbox.SetCell(x+w-1, y, sw.getDrawSymbol('┐'), sw.defaultForegroundColor, sw.defaultForegroundColor)
	termbox.SetCell(x+w-1, y+h-1, sw.getDrawSymbol('┘'), sw.defaultForegroundColor, sw.defaultForegroundColor)
}

func (sw *selectionWindow) drawSource(x, y, w int, data dto.GetDataResult) int {
	// draw caption aka source
	tbprint(x, y, sw.sourceForegroundColor, sw.sourceBackgroundColor, "Source: "+data.SourceName)
	y++

	// draw items
	for _, item := range data.Items {
		tbprint(x+1, y, sw.srKeyForegroundColor, sw.srKeyBackgroundColor, string(item.ButtonRune))
		tbprint(x+2, y, sw.defaultForegroundColor, sw.defaultBackgroundColor, " : "+item.Name)
		y++
	}
	return y
}

func (sw *selectionWindow) getDrawSymbol(s rune) rune {
	return sw.symbolsMap[s]
}

func (sw *selectionWindow) RemoveLastInput() {
	if !sw.focused || len(sw.currentInput) == 0 {
		return
	}
	sw.currentInput = sw.currentInput[:len(sw.currentInput)-1]
}

type SearchResultIface interface {
	GetSourceName() string
	GetCommands() []dto.CommandIface
	Founded() int
}
