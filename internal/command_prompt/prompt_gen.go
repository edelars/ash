package command_prompt

import (
	"encoding/json"
	"errors"
	"strings"

	"ash/internal/dto"

	"ash/pkg/termbox"

	"github.com/go-playground/colors"
)

const constErrParse string = "error parse config>"

var errEmptyValue = errors.New("empty value")

type promptItem struct {
	Value     string `json:"value"`
	Color     string `json:"color"` //#ffffff
	Bold      bool   `json:"bold"`
	Underline bool   `json:"underline"`
}

func parsePromptConfigString(b []byte) (res []promptItem) {
	if err := json.Unmarshal(b, &res); err != nil {
		res = append(res, promptItem{Value: constErrParse})
	}
	return res
}

func (c *CommandPrompt) generatePrompt(iContext dto.InternalContextIface) (res []termbox.Cell) {
	for _, item := range c.template {
		c, err := c.generatePieceOfPrompt(iContext, item)
		res = append(res, c...)
		if err != nil {
			break
		}
	}
	return res
}

func (c *CommandPrompt) generatePieceOfPrompt(iContext dto.InternalContextIface, item promptItem) ([]termbox.Cell, error) {
	var color colors.Color
	color, err := colors.ParseHEX(item.Color)
	if err != nil {
		color = nil
	}

	switch item.Value {
	case "":
		return c.stringToCells("empty value", nil, false, false), errEmptyValue
	case hasPrefix(item.Value, "%"):
		exeRes, err := c.execAdapter.ExecCmd(iContext, item.Value)
		return c.stringToCells(exeRes, color, item.Bold, item.Underline), err
	case hasPrefix(item.Value, "$"):
		return c.stringToCells(iContext.GetVariable(dto.Variable(item.Value)), color, item.Bold, item.Underline), nil
	default:
		return c.stringToCells(item.Value, color, item.Bold, item.Underline), nil

	}
}

func (c *CommandPrompt) stringToCells(s string, color colors.Color, b, u bool) (res []termbox.Cell) {
	var fg termbox.Attribute

	if color == nil {
		fg = c.defaultForegroundColor
	} else {
		fg = termbox.RGBToAttribute(color.ToRGB().R, color.ToRGB().G, color.ToRGB().B)
	}

	if b {
		fg = fg | termbox.AttrBold
	}
	if u {
		fg = fg | termbox.AttrUnderline
	}
	for _, r := range s {
		res = append(res, termbox.Cell{Ch: r, Fg: fg, Bg: c.defaultBackgroundColor})
	}
	return res
}

func hasPrefix(s string, p string) string {
	if strings.HasPrefix(s, p) {
		return s
	}
	return ""
}

type executionAdapter interface {
	ExecCmd(iContext dto.InternalContextIface, cmd string) (string, error)
}
