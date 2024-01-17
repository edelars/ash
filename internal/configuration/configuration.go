package configuration

import (
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
	"gopkg.in/yaml.v3"
)

const (
	constMainConfigDefaultFilename = "ash"
	constMainConfigDefaultDir      = "ash"
)

type ConfigLoader struct {
	Keybindings []KeyBind `yaml:"keybindings"`
	Aliases     []Alias   `yaml:"aliases"`
}

type KeyBind struct {
	Key    int    `yaml:"key"`
	Action string `yaml:"action"`
}

type Alias struct {
	Short string `yaml:"short"`
	Full  string `yaml:"full"`
}

func NewConfigLoader(errs chan error) ConfigLoader {
	startupConfig := newStartupConfigLoader()

	mainConfigFilename := getConfigFilename(startupConfig.Options.ConfigDir, configdir.LocalConfig())
	var config ConfigLoader
	// Read the file
	data, err := os.ReadFile(mainConfigFilename)
	if err != nil {
		errs <- err
		return config
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		errs <- err
		return config
	}
	// fmt.Println(config)
	return config
}

func getConfigFilename(startupFilename string, defaultConfigDir string) string {
	if startupFilename == "" {
		startupFilename = filepath.Join(defaultConfigDir, constMainConfigDefaultDir, constMainConfigDefaultFilename)
	}

	return startupFilename
}
