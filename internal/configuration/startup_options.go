package configuration

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	ConfigDir string `short:"c" long:"cfg" description:"Main configuration file" default:""`
}

type startupConfigLoader struct {
	Options
}

func newStartupConfigLoader() startupConfigLoader {
	var env Options

	p := flags.NewParser(&env, flags.Default)
	if _, err := p.Parse(); err != nil {
		fmt.Print("%w", err)
	}
	return startupConfigLoader{env}
}
