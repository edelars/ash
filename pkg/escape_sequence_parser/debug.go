package escape_sequence_parser

import (
	"ash/internal/configuration"
	"encoding/hex"
	"fmt"
	"os"
)

type esDebug struct {
	next EscapeSequenceParserIface
	file *os.File
}

func (e *esDebug) ParseEscapeSequence(buf []byte) []EscapeSequenceResultIface {
	res := e.next.ParseEscapeSequence(buf)

	if e.file != nil && len(buf) > 1 && buf[0] == 0x1b {
		e.file.WriteString(fmt.Sprintf("! block begin, raw binary data: %v\n", hex.EncodeToString(buf)))
		e.file.WriteString(fmt.Sprintf("count actions: %d \n", len(res)))
		for c, v := range res {
			a := string(rune(v.GetAction()))
			if v.GetAction() == 0x00 {
				a = "none"
			}
			e.file.WriteString(fmt.Sprintf("#%d action type: %s \n", c+1, a))
			e.file.WriteString(fmt.Sprintf("#%d raw args: %s \n", c+1, string(v.GetRaw())))
			n1, n2 := v.GetIntsFromArgs()
			e.file.WriteString(fmt.Sprintf("#%d int args: %d, %d \n\n", c+1, n1, n2))
		}
		e.file.WriteString(fmt.Sprintf("! end block\n\n"))
	}
	return res
}

func (e *esDebug) Stop() {
	if e.file != nil {
		e.file.Close()
	}
}

func NewESDebug(next EscapeSequenceParserIface, debugOpts configuration.Debug) esDebug {
	r := esDebug{next: next, file: nil}
	if debugOpts.EscapeSequence && debugOpts.DebugLogFile != "" {
		file, err := os.Create(debugOpts.DebugLogFile)
		if err != nil {
			fmt.Printf("Unable to create debug log file:%s %s\n", debugOpts.DebugLogFile, err)
			os.Exit(1)
		}
		r.file = file
	}
	return r
}
