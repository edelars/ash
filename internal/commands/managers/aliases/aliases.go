package aliases

import (
	"fmt"

	"ash/internal/commands"
	"ash/internal/dto"
)

const (
	constManagerName = "Aliases"
)

func NewAliasesManager(aDump aliasDump, executor executorIface, executeKey uint16, iContext dto.InternalContextIface) commands.CommandManagerIface {
	var cmds []dto.CommandIface

	for _, alias := range aDump.GetAliases() {
		cmds = append(cmds, NewAliasCommand(alias.Full, generateDescription(alias.Full), alias.Short, executor, executeKey, iContext))
	}

	return commands.NewCommandManager(constManagerName, 1, false, cmds...)
}

type aliasDump interface {
	GetAliases() []struct {
		Short string
		Full  string
	}
}

func generateDescription(full string) string {
	return fmt.Sprintf("alias for: %s", full)
}
