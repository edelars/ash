package commands_test

import (
	"testing"

	"ash/internal/commands"
	"ash/internal/dto"

	"github.com/stretchr/testify/assert"
)

func TestCommandRouterSearchResult_addResult(t *testing.T) {
	r := commands.NewCommandRouterSearchResult()
	p1 := commands.NewPattern("qweq", false)
	c := CommandImpl{W: 44}
	s1 := ManagerSearchResultImpl{
		Source:   "3333",
		Commands: []dto.CommandIface{&c},
		Pattern:  p1,
	}
	r.AddResult(s1)

	assert.Equal(t, 1, len(r.GetDataByPattern(p1)))
	cc := r.GetDataByPattern(p1)
	assert.Equal(t, 1, len(cc))
	assert.Equal(t, int8(44), cc[0].GetCommands()[0].GetMathWeight())

	p2 := commands.NewPattern("qw123123", false)
	cc2 := r.GetDataByPattern(p2)
	assert.Equal(t, 0, len(cc2))
}

type ManagerSearchResultImpl struct {
	Source   string
	Commands []dto.CommandIface
	Pattern  dto.PatternIface
}

func (managersearchresult ManagerSearchResultImpl) GetSourceName() string {
	return managersearchresult.Source
}

func (managersearchresult ManagerSearchResultImpl) GetCommands() []dto.CommandIface {
	return managersearchresult.Commands
}

func (managersearchresult ManagerSearchResultImpl) GetPattern() dto.PatternIface {
	return managersearchresult.Pattern
}

func (managersearchresult ManagerSearchResultImpl) Founded() int {
	return len(managersearchresult.Commands)
}

type CommandImpl struct {
	W int8
}

func (commandimpl *CommandImpl) SetMathWeight(weight int8) {
	panic("not implemented") // TODO: Implement
}

func (commandimpl *CommandImpl) GetName() string {
	panic("not implemented") // TODO: Implement
}

func (commandimpl *CommandImpl) GetExecFunc() dto.ExecF {
	panic("not implemented") // TODO: Implement
}

func (commandimpl CommandImpl) GetMathWeight() int8 {
	return commandimpl.W
}

func (commandimpl CommandImpl) WithArgs(args []string) dto.CommandIface {
	panic("not implemented")
}

func (commandimpl CommandImpl) GetArgs() []string {
	panic("not implemented")
}

func (commandimpl CommandImpl) MustPrepareExecutionList() bool {
	panic("not implemented")
}

type CommandManagerImpl struct{}

func (commandmanagerimpl CommandManagerImpl) SearchCommands(iContext dto.InternalContextIface, resultChan chan dto.CommandManagerSearchResult, patterns ...dto.PatternIface) {
	c := CommandImpl{W: 44}
	resultChan <- ManagerSearchResultImpl{
		Source:   "2",
		Commands: []dto.CommandIface{&c},
		Pattern:  patterns[0],
	}
}

type CommandManagerImpl2 struct{}

func (commandmanagerimpl CommandManagerImpl2) SearchCommands(iContext dto.InternalContextIface, resultChan chan dto.CommandManagerSearchResult, patterns ...dto.PatternIface) {
	resultChan <- ManagerSearchResultImpl{
		Source:   "3",
		Commands: []dto.CommandIface{},
		Pattern:  patterns[0],
	}
}

func TestCommandRouter_SearchCommands(t *testing.T) {
	cr := commands.NewCommandRouter(CommandManagerImpl{}, CommandManagerImpl2{})
	p1 := commands.NewPattern("123", false)

	r := cr.SearchCommands(nil, p1)
	cc := r.GetDataByPattern(p1)

	for _, v := range cc {
		assert.Equal(t, 1, v.Founded())
	}
	assert.Equal(t, 1, len(cc))
	assert.Equal(t, p1, cc[0].GetPattern())
	assert.Equal(t, "2", cc[0].GetSourceName())
	assert.Equal(t, int8(44), cc[0].GetCommands()[0].GetMathWeight())
}
