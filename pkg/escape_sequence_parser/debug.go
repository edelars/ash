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

func (esdebug *esDebug) ParseEscapeSequence(buf []byte) []EscapeSequenceResultIface {
	res := esdebug.next.ParseEscapeSequence(buf)

	if esdebug.file != nil && len(buf) > 1 && buf[0] == 0x1b {
		esdebug.file.WriteString(fmt.Sprintf("data: %v\n", hex.EncodeToString(buf)))
		esdebug.file.WriteString(fmt.Sprintf("count actions: %d \n", len(res)))
		for c, v := range res {
			a := string(rune(v.GetAction()))
			if v.GetAction() == 0x00 {
				a = "none"
			}
			esdebug.file.WriteString(fmt.Sprintf("N %d action type: %s \n", c+1, a))
			// fmt.Sprintf("N %d action type: %s \n", c+1, strconv.FormatInt(int64(v.GetAction()), 16)),

			esdebug.file.WriteString(fmt.Sprintf("N %d raw args: %s \n", c+1, string(v.GetRaw())))
			n1, n2 := v.GetIntsFromArgs()
			esdebug.file.WriteString(fmt.Sprintf("N %d int args: %d, %d \n", c+1, n1, n2))
			esdebug.file.WriteString(fmt.Sprintf("\n"))
		}
	}
	return res
}

func (esdebug *esDebug) Stop() {
	esdebug.file.Close()
}

func NewESDebug(next EscapeSequenceParserIface, debugOpts configuration.Debug) esDebug {
	r := esDebug{next: next, file: nil}
	if debugOpts.EscapeSequence && debugOpts.DebugLogFile != "" {
		file, err := os.Create(debugOpts.DebugLogFile)
		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		r.file = file
	}
	return r
}
