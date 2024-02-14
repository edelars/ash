package commands

func NewPattern(pattern string, precisionSearch bool) Pattern {
	return Pattern{
		value:           pattern,
		precisionSearch: precisionSearch,
	}
}

type Pattern struct {
	value           string
	precisionSearch bool
}

func (p Pattern) GetPattern() string {
	return p.value
}

func (p Pattern) IsPrecisionSearch() bool {
	if p.value == "" {
		return false
	}
	return p.precisionSearch
}
