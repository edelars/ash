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

const (
	CmdExecute          = ":Execute"
	CmdClose            = ":Close"
	CmdAutocomplete     = ":Autocomplete"
	CmdRemoveLeftSymbol = ":RemoveLeftSymbol"

	CmdKeyUp    = ":ArrowKeyUp"
	CmdKeyDown  = ":ArrowKeyDown"
	CmdKeyLeft  = ":ArrowKeyLeft"
	CmdKeyRight = ":ArrowKeyRight"
)

type ConfigLoader struct {
	ConfigFileName string            // for 'config' command output
	Prompt         string            `yaml:"prompt"`
	Keybindings    []KeyBind         `yaml:"keybindings"`
	Aliases        []Alias           `yaml:"aliases"`
	Envs           []string          `yaml:"envs"`
	Colors         Colors            `yaml:"colors"`
	Autocomplete   AutocompleteOpts  `yaml:"autocomplete"`
	Sqlite         StorageSqliteOpts `yaml:"sqlite"`
}

type Colors struct {
	DefaultText       string `yaml:"defaultText"`
	DefaultBackground string `yaml:"defaultBackground"`

	AutocompleteColors AutocompleteColors `yaml:"autocomplete"`
}

type AutocompleteColors struct {
	SourceText       string `yaml:"sourceText"`
	SourceBackground string `yaml:"sourceBackground"`
	ResultKeyText    string `yaml:"resultKeyText"`
	ResultBackground string `yaml:"resultBackground"`
}

type AutocompleteOpts struct {
	ShowFileInformation   bool `yaml:"showFileInformation"`
	InputFocusedByDefault bool `yaml:"inputFocusedByDefault"`
	ColumnGap             int  `yaml:"columnGap"`
}

type StorageSqliteOpts struct {
	FileName         string `yaml:"file"`
	WriteBuffer      int    `yaml:"writeBuffer"`
	MaxHistoryPerDir int    `yaml:"maxHistoryPerDir"`
	MaxHistoryTotal  int    `yaml:"maxHistoryTotal"`
	CleanupInterval  int    `yaml:"cleanupInterval"`
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
	config := newConfigLoaderWithDefaults()

	if _, err := os.Stat(mainConfigFilename); err != nil {
		return config
	}
	data, err := os.ReadFile(mainConfigFilename)
	if err != nil {
		return config
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config
	}

	config.ConfigFileName = mainConfigFilename
	return config
}

func newConfigLoaderWithDefaults() ConfigLoader {
	c := ConfigLoader{
		Keybindings: []KeyBind{
			{27, ":Close"},
			{13, ":Execute"},
			{9, ":Autocomplete"},
			{127, ":RemoveLeftSymbol"},
			{65514, ":ArrowKeyRight"},
			{65515, ":ArrowKeyLeft"},
			{65516, ":ArrowKeyDown"},
			{65517, ":ArrowKeyUp"},
		},
		Prompt: "ASH- ",
		Colors: Colors{
			DefaultText: "none", DefaultBackground: "none",
			AutocompleteColors: AutocompleteColors{
				SourceText:       "none",
				SourceBackground: "#8ec07c",
				ResultKeyText:    "none",
				ResultBackground: "#fabd2f",
			},
		},
		Autocomplete: AutocompleteOpts{
			ShowFileInformation: true, InputFocusedByDefault: false, ColumnGap: 3,
		},
		Sqlite: StorageSqliteOpts{
			FileName:         "sqlite.db",
			WriteBuffer:      3,
			MaxHistoryPerDir: 10,
			MaxHistoryTotal:  1000,
			CleanupInterval:  60,
		},
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
