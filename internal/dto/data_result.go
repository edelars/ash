package dto

type DataSource interface {
	GetCommand(r rune) CommandIface
	GetData(avalaibleSpace, overheadLinesPerSource int) []GetDataResult
}

type GetDataResult struct {
	SourceName string
	Items      []GetDataResultItem
}

type GetDataResultItem struct {
	Name       string
	ButtonRune rune
}
