package command_prompt

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"ash/internal/dto"
)

func newExecAdapter(exec Executor, key uint16) *execAdapter {
	r := &execAdapter{exec: exec, key: key}
	return r
}

type execAdapter struct {
	exec Executor
	key  uint16
}

// Adapter for exec commands like %(ls -la) via internal execution system
func (a *execAdapter) ExecCmd(iContext dto.InternalContextIface, c string) (res string, err error) {
	var cmd string
	if cmd, err = extractCmd(c); err != nil {
		return "", err
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	ictx := iContext.WithLastKeyPressed(a.key).WithCurrentInputBuffer([]rune(cmd)).WithOutputWriter(w)
	r := a.exec.Execute(ictx)

	if r != dto.CommandExecResultStatusOk {
		return "", errors.New(fmt.Sprintf("exit status %d", r))
	}

	if err = w.Flush(); err != nil {
		return "", err
	}

	return strings.ReplaceAll(b.String(), "\n", ""), nil
}

func extractCmd(s string) (string, error) {
	if s, finded := strings.CutPrefix(s, "%("); finded {
		if s, finded = strings.CutSuffix(s, ")"); finded {
			return s, nil
		}
	}
	return "", errors.New("no cmd pattern %() found")
}
