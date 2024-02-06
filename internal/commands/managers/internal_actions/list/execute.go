package list

import (
	"bytes"
	"io"
	"sync"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewExecuteCommand() *commands.Command {
	return commands.NewCommand(":Execute",
		func(iContext dto.InternalContextIface, _ []string) int {
			return executeCommands(iContext, nil, execCmd)
		}, true)
}

func execCmd(iContext dto.InternalContextIface, r io.Reader, w *io.PipeWriter, cmd dto.CommandIface) (st stResult) {
	var b []byte
	secondWriter := bytes.NewBuffer(b)
	mWriter := io.MultiWriter(secondWriter, w)
	newIC := iContext.WithInputReader(r).WithOutputWriter(mWriter)
	st.code = cmd.GetExecFunc()(newIC, cmd.GetArgs())
	if st.code > 0 {
		st.output = secondWriter.Bytes()
	}
	w.Close()

	return
}

type stResult struct {
	code   int
	output []byte
}

func executeCommands(iContext dto.InternalContextIface, _ []string, execFunc func(iContext dto.InternalContextIface, r io.Reader, w *io.PipeWriter, cmd dto.CommandIface) stResult) int {
	var lastReaderPipe *io.PipeReader
	var wg sync.WaitGroup
	var readerPipe1, readerPipe2 *io.PipeReader
	var writerPipe1, writerPipe2 *io.PipeWriter
	readerPipe1, writerPipe1 = io.Pipe()

	lastReaderPipe = readerPipe1

	b := true

	doneChan := make(chan struct{}, 1)

	returnChan := make(chan struct{})

	cmds := iContext.GetExecutionList()
	resChan := make(chan stResult, len(cmds))

	wg.Add(len(cmds))
	go func() {
		wg.Wait()
		returnChan <- struct{}{}
	}()

	for i := 0; i < len(cmds); i++ {
		if i == len(cmds)-1 { // last one
			doneChan <- struct{}{}
		}
		if i == 0 {
			go func(it int) {
				resChan <- execFunc(iContext, iContext.GetInputReader(), writerPipe1, cmds[it])
				wg.Done()
			}(i)
		} else {
			if b {
				readerPipe2, writerPipe2 = io.Pipe()
				lastReaderPipe = readerPipe2

				go func(it int) {
					resChan <- execFunc(iContext, readerPipe1, writerPipe2, cmds[it])
					wg.Done()
				}(i)
			} else {
				readerPipe1, writerPipe1 = io.Pipe()
				lastReaderPipe = readerPipe1

				go func(it int) {
					resChan <- execFunc(iContext, readerPipe2, writerPipe1, cmds[it])
					wg.Done()
				}(i)
			}
			b = !b
		}
	}

	var res int

	var doneBuf []byte
	defer func() {
		if res == 0 {
			iContext.GetOutputWriter().Write(doneBuf)
		}
		close(resChan)
		close(returnChan)
		close(doneChan)
	}()

	for {
		select {
		case st := <-resChan:
			if st.code > 0 {
				iContext.GetOutputWriter().Write(st.output)
				res = st.code
			}
		case <-returnChan:
			return res
		case <-doneChan:
			buf := make([]byte, 1024)
			for {
				n, err := lastReaderPipe.Read(buf)
				if err == io.EOF {
					break
				}
				doneBuf = buf[:n]
			}

		}
	}
}
