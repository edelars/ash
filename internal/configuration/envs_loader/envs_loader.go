package envs_loader

import (
	"errors"
	"os"
	"strings"
)

func LoadEnvs(dumper envsDumper) {
	for _, envStr := range dumper.GetEnvs() {
		if a, b, err := ParseEnvString(envStr); err == nil {
			os.Setenv(a, b)
		}
	}
}

type envsDumper interface {
	GetEnvs() []string
}

func ParseEnvString(s string) (a, b string, err error) {
	var found bool
	a, b, found = strings.Cut(strings.TrimSpace(s), "=")

	if !found {
		return "", "", errors.New("fail")
	}
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)

	return
}
