package commands

func NewEmptyCommand(r, description, displayName string) *Command {
	return NewCommandWithExtendedInfo(r, nil, false, description, displayName)
}
