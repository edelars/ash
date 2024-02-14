package list

import (
	"bytes"
	"io"
	"sync"

	"ash/internal/commands"
	"ash/internal/dto"
	"ash/internal/storage"

	"github.com/zenthangplus/goccm"
)

func NewExecuteCommand(historyAddFunc func(data storage.DataIface)) *commands.Command {
	return commands.NewCommand(":Execute",
		func(iContext dto.InternalContextIface, _ []string) dto.ExecResult {
			historyAddFunc(iContext)
			return executeCommands(iContext, nil, execCmd)
		}, true)
}

func execCmd(iContext dto.InternalContextIface, r io.Reader, w *io.PipeWriter, cmd dto.CommandIface) (st stResult) {
	b := make([]byte, 1024)
	secondWriter := bytes.NewBuffer(b)
	mWriter := io.MultiWriter(secondWriter, w)
	newIC := iContext.WithInputReader(r).WithOutputWriter(mWriter)
	st.code = cmd.GetExecFunc()(newIC, cmd.GetArgs())
	if st.code != dto.CommandExecResultStatusOk {
		st.output = secondWriter.Bytes()
	}
	w.Close()

	return
}

type stResult struct {
	code   dto.ExecResult
	output []byte
}

func executeCommands(iContext dto.InternalContextIface, _ []string, execFunc func(iContext dto.InternalContextIface, r io.Reader, w *io.PipeWriter, cmd dto.CommandIface) stResult) dto.ExecResult {
	cmds := iContext.GetExecutionList()
	if len(cmds) == 0 {
		return dto.CommandExecResultStatusOk
	}
	var lastReaderPipe *io.PipeReader

	ccm := goccm.New(2)
	var wg sync.WaitGroup
	var readerPipe1, readerPipe2 *io.PipeReader
	var writerPipe1, writerPipe2 *io.PipeWriter
	readerPipe1, writerPipe1 = io.Pipe()

	lastReaderPipe = readerPipe1

	doneChan := make(chan struct{}, 1)

	returnChan := make(chan struct{})

	resChan := make(chan stResult, len(cmds))

	defer func() {
		close(resChan)
		close(returnChan)
		close(doneChan)
	}()

	b := true
	for i := 0; i < len(cmds); i++ {
		ccm.Wait()

		cmd := cmds[i]
		if i == 0 {
			go func() {
				resChan <- execFunc(iContext, iContext.GetInputReader(), writerPipe1, cmd)
				ccm.Done()
			}()
		} else {
			switch b {
			case true:
				readerPipe2, writerPipe2 = io.Pipe()
				lastReaderPipe = readerPipe2

				go func() {
					resChan <- execFunc(iContext, readerPipe1, writerPipe2, cmd)
					ccm.Done()
				}()
			case false:
				readerPipe1, writerPipe1 = io.Pipe()
				lastReaderPipe = readerPipe1

				go func() {
					resChan <- execFunc(iContext, readerPipe2, writerPipe1, cmd)
					ccm.Done()
				}()
			}
			b = !b
		}

		if i == len(cmds)-1 { // last one, starting read from last goroutine
			doneChan <- struct{}{}
		}

	}

	wg.Add(len(cmds))
	go func() {
		wg.Wait()
		returnChan <- struct{}{}
	}()

	res := dto.CommandExecResultStatusOk
	var doneBuf []byte

mainLoop:
	for {
		select {
		case st := <-resChan:
			if st.code != dto.CommandExecResultStatusOk && res == dto.CommandExecResultStatusOk {
				iContext.GetOutputWriter().Write(st.output)
				res = st.code
			}
			wg.Done()
		case <-returnChan:
			break mainLoop
		case <-doneChan:
			buf := make([]byte, 1024)
		readLoop:
			for {
				n, err := lastReaderPipe.Read(buf)
				if err == io.EOF {
					break readLoop
				}
				doneBuf = append(doneBuf, buf[:n]...)
			}
		}
	}

	if res == dto.CommandExecResultStatusOk {
		iContext.GetOutputWriter().Write(doneBuf)
	}
	return res
}

type storageIface interface {
	SaveData(data dataIface)
}

type dataIface interface {
	GetExecutionList() []dto.CommandIface
	GetCurrentDir() string
}
