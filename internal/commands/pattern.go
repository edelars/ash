package commands

type PatternIface interface {
	GetPattern() string
}

func NewPattern(pattern string) Pattern {
	return Pattern{
		value: pattern,
	}
}

type Pattern struct {
	value string
}

func (p Pattern) GetPattern() string {
	return p.value
}
