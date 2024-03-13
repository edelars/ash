package history

import (
	"io"
	"testing"

	"ash/internal/commands"
	"ash/internal/dto"
	"ash/internal/storage"

	"ash/pkg/termbox"

	"github.com/stretchr/testify/assert"
)

func Test_historyManager_convertStorageItems(t *testing.T) {
	h := historyManager{}
	res := h.convertStorageItems([]storage.StorageResult{&storageResImpl{"cmd1", "dir1", 2}, &storageResImpl{"cmd2", "dir2", 3}})

	assert.Equal(t, 2, len(res))
	assert.Equal(t, "cmd1", res[0].GetName())
	assert.Equal(t, "cmd1", res[0].GetDisplayName())
	assert.Equal(t, uint8(0), res[0].GetMathWeight())

	assert.Equal(t, generateDescription("dir1", 2), res[0].GetDescription())

	assert.Equal(t, "cmd2", res[1].GetName())
	assert.Equal(t, "cmd2", res[1].GetDisplayName())
	assert.Equal(t, uint8(1), res[1].GetMathWeight())
	assert.Equal(t, generateDescription("dir2", 3), res[1].GetDescription())

	res = h.convertStorageItems(nil)
	assert.Equal(t, 0, len(res))
}

type storageResImpl struct {
	c, d string
	u    int
}

func (storageresimpl *storageResImpl) GetCommand() string {
	return storageresimpl.c
}

func (storageresimpl *storageResImpl) GetDir() string {
	return storageresimpl.d
}

func (storageresimpl *storageResImpl) GetUsedCount() int {
	return storageresimpl.u
}

func Test_historyManager_SearchCommands(t *testing.T) {
	ch := make(chan dto.CommandManagerSearchResult)
	defer close(ch)

	s := storageImpl{}
	h := historyManager{}
	h.storage = &s
	go h.SearchCommands(&contImpl{}, ch, commands.NewPattern("f1", false))
	<-ch
	assert.Equal(t, false, s.getTopDirs)
	assert.Equal(t, true, s.getTopPattern)

	s.getTopPattern = false
	s.getTopDirs = false
	go h.SearchCommands(&contImpl{}, ch, commands.NewPattern("", false))
	<-ch
	assert.Equal(t, true, s.getTopDirs)
	assert.Equal(t, false, s.getTopPattern)

	s.getTopPattern = false
	s.getTopDirs = false
	go h.SearchCommands(&contImpl{}, ch, commands.NewPattern("asdf", false))
	<-ch
	assert.Equal(t, false, s.getTopDirs)
	assert.Equal(t, true, s.getTopPattern)
}

type storageImpl struct {
	getTopDirs, getTopPattern bool
}

func (storageimpl *storageImpl) SaveData(data storage.DataIface) {
	panic("not implemented") // TODO: Implement
}

func (storageimpl *storageImpl) GetTopHistoryByDirs(currentDir string, l int) []storage.StorageResult {
	storageimpl.getTopDirs = true
	return nil
}

func (storageimpl *storageImpl) GetTopHistoryByPattern(prefix string, l int) []storage.StorageResult {
	storageimpl.getTopPattern = true
	return nil
}

type contImpl struct{}

// ctrl-c - break app
func (contimpl *contImpl) GetExecTerminateChan() chan struct{} {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) WithVariables(vars []dto.VariableSet) dto.InternalContextIface {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) GetCellsPrintFunction() func(cells []termbox.Cell) {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) GetVariable(v dto.Variable) string {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) GetEnv(envName string) string {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) GetCurrentDir() string {
	return "123"
}

func (contimpl *contImpl) WithLastKeyPressed(b uint16) dto.InternalContextIface {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) WithCurrentInputBuffer(b []rune) dto.InternalContextIface {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) GetCurrentInputBuffer() []rune {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) GetLastKeyPressed() uint16 {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) GetInputEventChan() chan termbox.Event {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) GetErrChan() chan error {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) WithExecutionList(executionList []dto.CommandIface) dto.InternalContextIface {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) GetExecutionList() []dto.CommandIface {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) GetPrintFunction() func(msg string) {
	panic("not implemented") // TODO: Implement
}

// console I/O
func (contimpl *contImpl) GetOutputWriter() io.Writer {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) GetInputReader() io.Reader {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) WithOutputWriter(_ io.Writer) dto.InternalContextIface {
	panic("not implemented") // TODO: Implement
}

func (contimpl *contImpl) WithInputReader(_ io.Reader) dto.InternalContextIface {
	panic("not implemented") // TODO: Implement
}
