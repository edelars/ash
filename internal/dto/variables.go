package dto

type Variable string

type VariableSet struct {
	Name  Variable
	Value string
}

const (
	VariableCurDir       Variable = "$CURDIR"
	VariableLastExitCode Variable = "$?"
	VariableCurrentUser  Variable = "$USER"
	VariableHostname     Variable = "$HOSTNAME"
)
