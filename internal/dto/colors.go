package dto

import "ash/pkg/termbox"

type ColorsAdapterIface interface {
	GetColors() Colors
}

type Colors struct {
	DefaultForegroundColor  termbox.Attribute
	DefaultBackgroundColor  termbox.Attribute
	SelectedForegroundColor termbox.Attribute
	AutocompleteColors      AutocompleteColors
}

type AutocompleteColors struct {
	SourceBackgroundColor    termbox.Attribute
	SourceForegroundColor    termbox.Attribute
	ResultKeyBackgroundColor termbox.Attribute
	ResultKeyForegroundColor termbox.Attribute
	DescriptionText          termbox.Attribute
}
