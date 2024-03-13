package dto

type DataSource interface {
	GetCommand(r rune) CommandIface
	GetData(avalaibleSpace, overheadLinesPerSource int) ([]GetDataResult, int, int) // data, mainFieldMaxWid,descFieldMaxWid
}

type GetDataResult struct {
	SourceName string
	Items      []GetDataResultItem
}

type GetDataResultItem struct {
	ButtonRune  rune
	Name        string
	DisplayName string
	Description string
}
