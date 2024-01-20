package configuration

import (
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
	"gopkg.in/yaml.v3"
)

const (
	constMainConfigDefaultFilename = "ash.yaml"
	constMainConfigDefaultDir      = "ash"
)

type ConfigLoader struct {
	Keybindings []KeyBind `yaml:"keybindings"`
	Aliases     []Alias   `yaml:"aliases"`
	Prompt      string    `yaml:"prompt"`
}

type KeyBind struct {
	Key    int    `yaml:"key"`
	Action string `yaml:"action"`
}

type Alias struct {
	Short string `yaml:"short"`
	Full  string `yaml:"full"`
}

func NewConfigLoader() ConfigLoader {
	startupConfig := newStartupConfigLoader()

	mainConfigFilename := getConfigFilename(startupConfig.Options.ConfigDir, configdir.LocalConfig())
	var config ConfigLoader

	if _, err := os.Stat(mainConfigFilename); err != nil {
		return newConfigLoaderWithDefaults()
	}
	// Read the file
	data, err := os.ReadFile(mainConfigFilename)
	if err != nil {
		return newConfigLoaderWithDefaults()
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return newConfigLoaderWithDefaults()
	}
	return config
}

func newConfigLoaderWithDefaults() ConfigLoader {
	c := ConfigLoader{
		Keybindings: []KeyBind{{13, ":Execute"}, {14, ":Autocomplete"}},
		Aliases:     []Alias{},
		Prompt:      "ASH> ",
	}
	return c
}

func getConfigFilename(startupFilename string, defaultConfigDir string) string {
	if startupFilename == "" {
		startupFilename = filepath.Join(defaultConfigDir, constMainConfigDefaultDir, constMainConfigDefaultFilename)
	}

	return startupFilename
}

func (c ConfigLoader) GetKeysBindings() []struct {
	Key    int
	Action string
} {
	var res []struct {
		Key    int
		Action string
	}

	for _, kb := range c.Keybindings {
		res = append(res, struct {
			Key    int
			Action string
		}{kb.Key, kb.Action})
	}
	return res
}
