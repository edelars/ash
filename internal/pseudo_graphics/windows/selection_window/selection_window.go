package selection_window

import (
	"ash/internal/configuration"
	"ash/internal/dto"

	"ash/pkg/termbox"

	"github.com/mattn/go-runewidth"
)

const (
	constHelpMessage   = "Press <tab> to change focus and <[a-z]> to execute command from search results"
	constNoDataMessage = "No search results"

	constColumMainMinWid   = 20
	constColumnFileInfoWid = 9
)

type selectionWindow struct {
	cursorX int
	cursorY int
	mainX   int
	mainY   int
	mainW   int
	mainH   int

	columnDescriptionX      int
	columnDescriptionMaxWid int
	columnGap               int

	defaultBackgroundColor   termbox.Attribute
	defaultForegroundColor   termbox.Attribute
	sourceBackgroundColor    termbox.Attribute
	sourceForegroundColor    termbox.Attribute
	resultKeyBackgroundColor termbox.Attribute
	resultKeyForegroundColor termbox.Attribute
	selectedForegroundColor  termbox.Attribute
	descriptionText          termbox.Attribute

	focused            bool
	needToUpdateSource bool
	dataSources        dto.DataSource
	symbolsMap         map[rune]rune
	currentInput       []rune
	searchFunc         func(patter []rune) dto.DataSource
	resultFunc         func(cmd dto.CommandIface, userInput []rune)

	showCommandDescription bool
}

func NewSelectionWindow(userInput []rune, searchFunc func(patter []rune) dto.DataSource, resultFunc func(cmd dto.CommandIface, userInput []rune), autocomplOpts configuration.AutocompleteOpts, colorsAdapter dto.ColorsAdapterIface) selectionWindow {
	colors := colorsAdapter.GetColors()

	sw := selectionWindow{
		focused:                  autocomplOpts.InputFocusedByDefault,
		symbolsMap:               map[rune]rune{},
		currentInput:             userInput,
		searchFunc:               searchFunc,
		resultFunc:               resultFunc,
		showCommandDescription:   autocomplOpts.ShowFileInformation,
		columnGap:                autocomplOpts.ColumnGap,
		needToUpdateSource:       true,
		defaultBackgroundColor:   colors.DefaultBackgroundColor,
		defaultForegroundColor:   colors.DefaultForegroundColor,
		sourceBackgroundColor:    colors.AutocompleteColors.SourceBackgroundColor,
		sourceForegroundColor:    colors.AutocompleteColors.SourceForegroundColor,
		resultKeyBackgroundColor: colors.AutocompleteColors.ResultKeyBackgroundColor,
		resultKeyForegroundColor: colors.AutocompleteColors.ResultKeyForegroundColor,
		selectedForegroundColor:  colors.SelectedForegroundColor,
		descriptionText:          colors.AutocompleteColors.DescriptionText,
	}
	if runewidth.EastAsianWidth {
		sw.symbolsMap = map[rune]rune{'─': '-', '│': '|', '┌': '+', '└': '+', '┐': '+', '┘': '+'}
	} else {
		sw.symbolsMap = map[rune]rune{'─': '─', '│': '│', '┌': '┌', '└': '└', '┐': '┐', '┘': '┘'}
	}

	return sw
}

func (sw *selectionWindow) Close() {
	sw.resultFunc(nil, sw.currentInput)
}

// Trim to last word ie: "ls /usr" to "/usr"
func trimInput(input []rune) []rune {
	var lastSpacePos int
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] == 32 {
			lastSpacePos = i
			break
		}
	}
	if lastSpacePos > 0 && len(input) > lastSpacePos {
		return input[lastSpacePos+1:]
	}
	if len(input) == 1 && input[0] == 32 { // single space
		return []rune(nil)
	}
	return input
}

// Trim last word ie: "ls /usr" to "ls "
func trimInputSuffix(input []rune) []rune {
	var lastSpacePos int
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] == 32 {
			lastSpacePos = i
			break
		}
	}
	if lastSpacePos > 0 {
		return input[:lastSpacePos+1]
	}
	if len(input) > 0 && input[len(input)-1] != 32 { // add space
		input = append(input, 32)
	}
	return input
}

func (sw *selectionWindow) updateDataSource() {
	if sw.needToUpdateSource {
		sw.needToUpdateSource = false
		sw.dataSources = sw.searchFunc(trimInput(sw.currentInput))
	}
}

func (sw *selectionWindow) KeyInput(key rune) {
	if sw.focused {
		sw.currentInput = append(sw.currentInput, rune(key))
		sw.needToUpdateSource = true
	} else {
		sw.commandChoosed(key)
	}
}

func (sw *selectionWindow) commandChoosed(key rune) {
	cmd := sw.dataSources.GetCommand(key)
	switch cmd.GetExecFunc() { // just change userInput. No execute
	case nil:
		withoutSuffix := trimInputSuffix(sw.currentInput)
		sw.currentInput = append(withoutSuffix, []rune(cmd.GetDisplayName())...)
		sw.needToUpdateSource = true
	default:
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
		termbox.SetFg(sw.cursorX, sw.cursorY, sw.selectedForegroundColor)
		termbox.SetCursor(sw.cursorX, sw.cursorY)
	} else {
		termbox.HideCursor()
		termbox.SetFg(sw.cursorX, sw.cursorY, sw.defaultForegroundColor)
	}
}

func (sw *selectionWindow) calculateColumnsWidth(mainFieldMaxWid, descFieldMaxWid int) {
	if sw.showCommandDescription && sw.mainW > mainFieldMaxWid {
		if descFieldMaxWid < sw.mainW-mainFieldMaxWid-sw.columnGap {
			sw.columnDescriptionX = mainFieldMaxWid + sw.columnGap + 1
			sw.columnDescriptionMaxWid = descFieldMaxWid
		} else if (sw.mainW - mainFieldMaxWid - sw.columnGap) > 10 {
			sw.columnDescriptionX = sw.mainW - mainFieldMaxWid - sw.columnGap + 1
			sw.columnDescriptionMaxWid = sw.mainW - mainFieldMaxWid - sw.columnGap
		}
	}
}

func (sw *selectionWindow) reDraw() {
	const bottomInputBarH = 2
	const overheadSpaceForSource = 2

	sw.updateDataSource()

	dataResult, mainFieldMaxWid, descFieldMaxWid := sw.dataSources.GetData(sw.mainH-bottomInputBarH, overheadSpaceForSource)
	sw.calculateColumnsWidth(mainFieldMaxWid, descFieldMaxWid)

	switch len(dataResult) {
	case 0: // if no data we will draw empty results message
		x := sw.mainW/2 - len([]rune(constNoDataMessage))/2 // calculate center of screen
		y := (sw.mainH - bottomInputBarH) / 2
		tbPrint(x, y, sw.defaultForegroundColor, sw.defaultBackgroundColor, constNoDataMessage)
	default: // draw data
		stepY := sw.mainY
		for _, ds := range dataResult {
			stepY = sw.drawSource(sw.mainX+1, stepY, sw.mainW, ds)
		}
	}

	// draw help
	tbPrint(sw.mainX+1, sw.mainY+sw.mainH-4, sw.defaultForegroundColor, sw.defaultBackgroundColor, constHelpMessage)
	// draw input
	sw.drawRectangle(sw.mainX, sw.mainY+sw.mainH-3, sw.mainW, 3) // draw input box
	curX := sw.mainX + 1
	curY := sw.mainY + sw.mainH - 2

	tbPrint(curX, curY, sw.defaultForegroundColor, sw.defaultBackgroundColor, string(sw.currentInput))
	curX = curX + len(sw.currentInput)
	sw.setCursor(curX, curY)
}

func (sw *selectionWindow) Draw(x, y, w, h int, fg, bg termbox.Attribute) {
	sw.mainX = x
	sw.mainY = y
	sw.mainW = w
	sw.mainH = h

	sw.defaultForegroundColor = fg
	sw.defaultBackgroundColor = bg

	sw.reDraw()
}

func (sw *selectionWindow) fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, sw.defaultForegroundColor, sw.defaultBackgroundColor)
		}
	}
}

func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) int {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += 1
	}
	return x
}

func tbPrintWithHighlights(x, y int, fg, bg termbox.Attribute, msg string, pattern []rune) int {
	var patternCursor int
	for _, c := range msg {
		f := fg
		if len(pattern) > patternCursor && c == pattern[patternCursor] {
			f = f | termbox.AttrUnderline
			patternCursor++
		}
		termbox.SetCell(x, y, c, f, bg)
		x += 1
	}
	return x
}

func (sw *selectionWindow) drawRectangle(x, y, w, h int) {
	sw.fill(x+1, y, w-2, 1, termbox.Cell{Ch: sw.getDrawSymbol('─')})     // top line
	sw.fill(x+1, y+h-1, w-2, 1, termbox.Cell{Ch: sw.getDrawSymbol('─')}) // bottom line

	sw.fill(x, y+1, 1, h-2, termbox.Cell{Ch: sw.getDrawSymbol('│')})     // left line
	sw.fill(x+w-1, y+1, 1, h-2, termbox.Cell{Ch: sw.getDrawSymbol('│')}) // right line

	termbox.SetCell(x, y, sw.getDrawSymbol('┌'), sw.defaultForegroundColor, sw.defaultBackgroundColor)
	termbox.SetCell(x, y+h-1, sw.getDrawSymbol('└'), sw.defaultForegroundColor, sw.defaultBackgroundColor)
	termbox.SetCell(x+w-1, y, sw.getDrawSymbol('┐'), sw.defaultForegroundColor, sw.defaultBackgroundColor)
	termbox.SetCell(x+w-1, y+h-1, sw.getDrawSymbol('┘'), sw.defaultForegroundColor, sw.defaultBackgroundColor)
}

func (sw *selectionWindow) drawSource(x, y, w int, data dto.GetDataResult) int {
	// draw caption aka source
	tbPrint(x, y, sw.sourceForegroundColor, sw.sourceBackgroundColor, "Source: "+data.SourceName)
	y++

	// draw items
	for _, item := range data.Items {
		xNew := tbPrint(x+1, y, sw.resultKeyForegroundColor|termbox.AttrBold, sw.resultKeyBackgroundColor, string(item.ButtonRune))
		xNew = tbPrint(xNew, y, sw.defaultForegroundColor, sw.defaultBackgroundColor, " : ")
		xNew = tbPrintWithHighlights(xNew, y, sw.defaultForegroundColor, sw.defaultBackgroundColor, item.DisplayName, trimInput(sw.currentInput))
		if sw.showCommandDescription && sw.columnDescriptionX > 0 && sw.columnDescriptionMaxWid > 0 {
			desc := firstN(item.Description, sw.columnDescriptionMaxWid)
			xNew = tbPrint(sw.columnDescriptionX+5, y, sw.descriptionText, sw.defaultBackgroundColor, desc)
		}
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

func firstN(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}
