package commands_test

import (
	"testing"

	"ash/internal/commands"

	"github.com/stretchr/testify/assert"
)

func TestCommandRouterSearchResult_addResult(t *testing.T) {
	r := commands.NewCommandRouterSearchResult()
	p1 := commands.NewPattern("qweq")
	c := CommandImpl{W: 44}
	s1 := ManagerSearchResultImpl{
		Source:   "3333",
		Commands: []commands.CommandIface{&c},
		Pattern:  p1,
	}
	r.AddResult(s1)

	assert.Equal(t, 1, len(r.GetDataByPattern(p1)))
	cc := r.GetDataByPattern(p1)
	assert.Equal(t, 1, len(cc))
	assert.Equal(t, int8(44), cc[0].GetCommands()[0].GetMathWeight())

	p2 := commands.NewPattern("qw123123")
	cc2 := r.GetDataByPattern(p2)
	assert.Equal(t, 0, len(cc2))
}

type ManagerSearchResultImpl struct {
	Source   string
	Commands []commands.CommandIface
	Pattern  commands.PatternIface
}

func (managersearchresult ManagerSearchResultImpl) GetSourceName() string {
	return managersearchresult.Source
}

func (managersearchresult ManagerSearchResultImpl) GetCommands() []commands.CommandIface {
	return managersearchresult.Commands
}

func (managersearchresult ManagerSearchResultImpl) GetPattern() commands.PatternIface {
	return managersearchresult.Pattern
}

type CommandImpl struct {
	W int8
}

func (commandimpl *CommandImpl) GetExecFunc() commands.ExecF {
	panic("not implemented") // TODO: Implement
}

func (commandimpl CommandImpl) GetMathWeight() int8 {
	return commandimpl.W
}

type CommandManagerImpl struct{}

func (commandmanagerimpl CommandManagerImpl) SearchCommands(resultChan chan commands.CommandManagerSearchResult, patterns ...commands.PatternIface) {
	c := CommandImpl{W: 44}
	resultChan <- ManagerSearchResultImpl{
		Source:   "2",
		Commands: []commands.CommandIface{&c},
		Pattern:  patterns[0],
	}
}

func TestCommandRouter_SearchCommands(t *testing.T) {
	cm := CommandManagerImpl{}
	cr := commands.NewCommandRouter(cm)
	p1 := commands.NewPattern("123")

	r := cr.SearchCommands(p1)
	cc := r.GetDataByPattern(p1)

	assert.Equal(t, 1, len(cc))
	assert.Equal(t, p1, cc[0].GetPattern())
	assert.Equal(t, "2", cc[0].GetSourceName())
	assert.Equal(t, int8(44), cc[0].GetCommands()[0].GetMathWeight())
}
