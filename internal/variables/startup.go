package variables

import (
	"os"
	"os/user"

	"ash/internal/dto"
)

func GetVariables() (res []dto.VariableSet) {
	if r := getUsername(); r != nil {
		res = append(res, *r)
	}

	if r := getHostname(); r != nil {
		res = append(res, *r)
	}
	return res
}

func getUsername() *dto.VariableSet {
	user, err := user.Current()
	if err != nil {
		return nil
	}
	return &dto.VariableSet{
		Name:  dto.VariableCurrentUser,
		Value: user.Name,
	}
}

func getHostname() *dto.VariableSet {
	name, err := os.Hostname()
	if err != nil {
		return nil
	}
	return &dto.VariableSet{
		Name:  dto.VariableHostname,
		Value: name,
	}
}
