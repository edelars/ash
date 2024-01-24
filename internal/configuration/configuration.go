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
	ConfigFileName string    // for 'config' command output
	Prompt         string    `yaml:"prompt"`
	Keybindings    []KeyBind `yaml:"keybindings"`
	Aliases        []Alias   `yaml:"aliases"`
	Envs           []string  `yaml:"envs"`
	Colors         Colors    `yaml:"colors"`
}

type Colors struct {
	DefaultText       uint64 `yaml:"defaultText"`
	DefaultBackground uint64 `yaml:"defaultBackground"`
	Autocomplete
}

type Autocomplete struct {
	SourceText       uint64 `yaml:"sourceText"`
	SourceBackground uint64 `yaml:"sourceBackground"`
	ResultKeyText    uint64 `yaml:"resultKeyText"`
	ResultBackground uint64 `yaml:"resultBackground"`
}

type KeyBind struct {
	Key    uint16 `yaml:"key"`
	Action string `yaml:"action"`
}

type Alias struct {
	Short string `yaml:"short"`
	Full  string `yaml:"full"`
}

func (c ConfigLoader) GetKeyBind(action string) uint16 {
	for _, v := range c.Keybindings {
		if v.Action == action {
			return v.Key
		}
	}
	return 0
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
	config.ConfigFileName = mainConfigFilename
	return config
}

func newConfigLoaderWithDefaults() ConfigLoader {
	c := ConfigLoader{
		Keybindings: []KeyBind{{27, ":Close"}, {13, ":Execute"}, {9, ":Autocomplete"}, {127, ":Backspace"}},
		Prompt:      "ASH> ",
		Colors: Colors{DefaultText: 0, DefaultBackground: 0, Autocomplete: Autocomplete{
			SourceText:       1,
			SourceBackground: 13,
			ResultKeyText:    1,
			ResultBackground: 11,
		}},
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
	Key    uint16
	Action string
} {
	var res []struct {
		Key    uint16
		Action string
	}

	for _, kb := range c.Keybindings {
		res = append(res, struct {
			Key    uint16
			Action string
		}{kb.Key, kb.Action})
	}
	return res
}

func (c ConfigLoader) GetEnvs() []string {
	return c.Envs
}

func (c ConfigLoader) GetConfig() interface{} {
	return c
}
