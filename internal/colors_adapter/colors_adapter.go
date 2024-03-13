package colors_adapter

import (
	"ash/internal/configuration"
	"ash/internal/dto"

	"ash/pkg/termbox"

	"github.com/go-playground/colors"
)

type colorsAdapter struct {
	defaultBackgroundColor   termbox.Attribute
	defaultForegroundColor   termbox.Attribute
	sourceBackgroundColor    termbox.Attribute
	sourceForegroundColor    termbox.Attribute
	resultKeyBackgroundColor termbox.Attribute
	resultKeyForegroundColor termbox.Attribute
	selectedForegroundColor  termbox.Attribute
	descriptionText          termbox.Attribute
}

func (c colorsAdapter) GetColors() dto.Colors {
	return dto.Colors{
		DefaultForegroundColor:  c.defaultForegroundColor,
		DefaultBackgroundColor:  c.defaultBackgroundColor,
		SelectedForegroundColor: c.selectedForegroundColor,
		AutocompleteColors: dto.AutocompleteColors{
			SourceBackgroundColor:    c.sourceBackgroundColor,
			SourceForegroundColor:    c.sourceForegroundColor,
			ResultKeyBackgroundColor: c.resultKeyBackgroundColor,
			ResultKeyForegroundColor: c.resultKeyForegroundColor,
			DescriptionText:          c.descriptionText,
		},
	}
}

func NewColorsAdapter(colorsOpts configuration.Colors) colorsAdapter {
	c := colorsAdapter{}
	// TODO add tests
	defaultBackgroundColor, err := colors.ParseHEX(colorsOpts.DefaultBackground)
	if err == nil {
		c.defaultBackgroundColor = termbox.RGBToAttribute(defaultBackgroundColor.ToRGB().R, defaultBackgroundColor.ToRGB().G, defaultBackgroundColor.ToRGB().B)
	}

	defaultForegroundColor, err := colors.ParseHEX(colorsOpts.DefaultText)
	if err == nil {
		c.defaultForegroundColor = termbox.RGBToAttribute(defaultForegroundColor.ToRGB().R, defaultForegroundColor.ToRGB().G, defaultForegroundColor.ToRGB().B)
	}

	sourceBackgroundColor, err := colors.ParseHEX(colorsOpts.AutocompleteColors.SourceBackground)
	if err == nil {
		c.sourceBackgroundColor = termbox.RGBToAttribute(sourceBackgroundColor.ToRGB().R, sourceBackgroundColor.ToRGB().G, sourceBackgroundColor.ToRGB().B)
	}

	sourceForegroundColor, err := colors.ParseHEX(colorsOpts.AutocompleteColors.SourceText)
	if err == nil {
		c.sourceForegroundColor = termbox.RGBToAttribute(sourceForegroundColor.ToRGB().R, sourceForegroundColor.ToRGB().G, sourceForegroundColor.ToRGB().B)
	}

	resultKeyBackgroundColor, err := colors.ParseHEX(colorsOpts.AutocompleteColors.ResultBackground)
	if err == nil {
		c.resultKeyBackgroundColor = termbox.RGBToAttribute(resultKeyBackgroundColor.ToRGB().R, resultKeyBackgroundColor.ToRGB().G, resultKeyBackgroundColor.ToRGB().B)
	}

	resultKeyForegroundColor, err := colors.ParseHEX(colorsOpts.AutocompleteColors.ResultKeyText)
	if err == nil {
		c.resultKeyForegroundColor = termbox.RGBToAttribute(resultKeyForegroundColor.ToRGB().R, resultKeyForegroundColor.ToRGB().G, resultKeyForegroundColor.ToRGB().B)
	}

	selectedForegroundColor, err := colors.ParseHEX(colorsOpts.SelectedForegroundColor)
	if err == nil {
		c.selectedForegroundColor = termbox.RGBToAttribute(selectedForegroundColor.ToRGB().R, selectedForegroundColor.ToRGB().G, selectedForegroundColor.ToRGB().B)
	}

	descriptionTextColor, err := colors.ParseHEX(colorsOpts.AutocompleteColors.DescriptionText)
	if err == nil {
		c.descriptionText = termbox.RGBToAttribute(descriptionTextColor.ToRGB().R, descriptionTextColor.ToRGB().G, descriptionTextColor.ToRGB().B)
	} else {
		// panic(colorsOpts.AutocompleteColors.DescriptionText)
	}
	return c
}
