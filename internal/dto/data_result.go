package dto

type GetDataResult struct {
	SourceName string
	Items      []GetDataResultItem
}

type GetDataResultItem struct {
	Name       string
	ButtonRune rune
}
