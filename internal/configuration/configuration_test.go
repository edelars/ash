package configuration

import (
	"path/filepath"
	"reflect"
	"testing"
)

func Test_getConfigFilename(t *testing.T) {
	type args struct {
		startupFilename  string
		defaultConfigDir string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "none",
			args: args{
				startupFilename:  "",
				defaultConfigDir: "/1",
			},
			want: filepath.Join("/1/", constMainConfigDefaultDir, constMainConfigDefaultFilename),
		},
		{
			name: "cfg.yaml",
			args: args{
				startupFilename:  "cfg.yaml",
				defaultConfigDir: "",
			},
			want: "cfg.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConfigFilename(tt.args.startupFilename, tt.args.defaultConfigDir); got != tt.want {
				t.Errorf("getConfigFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newConfigLoaderWithDefaults(t *testing.T) {
	tests := []struct {
		name string
		want ConfigLoader
	}{
		{
			name: "defaults",
			want: ConfigLoader{
				Keybindings: []KeyBind{{27, ":Close"}, {13, ":Execute"}, {9, ":Autocomplete"}, {127, ":RemoveLeftSymbol"}},
				Prompt:      "ASH- ",
				Colors:      Colors{},
				Autocomplete: AutocompleteOpts{
					ShowFileInformation: true, InputFocusedByDefault: false, ColumnGap: 3, Colors: AutocompleteColors{
						SourceText:       1,
						SourceBackground: 13,
						ResultKeyText:    1,
						ResultBackground: 11,
					},
				},
				Sqlite: StorageSqliteOpts{
					FileName:         "sqlite.db",
					WriteBuffer:      3,
					MaxHistoryPerDir: 10,
					MaxHistoryTotal:  1000,
					CleanupInterval:  60,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newConfigLoaderWithDefaults(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newConfigLoaderWithDefaults() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigLoader_GetKeysBindings(t *testing.T) {
	type fields struct {
		Keybindings []KeyBind
		Aliases     []Alias
		Prompt      string
	}
	tests := []struct {
		name   string
		fields fields
		want   []struct {
			Key    uint16
			Action string
		}
	}{
		{
			name: "empty",
			fields: fields{
				Keybindings: []KeyBind{},
				Aliases:     []Alias{},
				Prompt:      "",
			},
			want: nil,
		},
		{
			name: "1",
			fields: fields{
				Keybindings: []KeyBind{{13, "enter"}, {11, "done"}},
				Aliases:     []Alias{},
				Prompt:      "",
			},
			want: []struct {
				Key    uint16
				Action string
			}{{13, "enter"}, {11, "done"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ConfigLoader{
				Keybindings: tt.fields.Keybindings,
				Aliases:     tt.fields.Aliases,
				Prompt:      tt.fields.Prompt,
			}
			if got := c.GetKeysBindings(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigLoader.GetAliases() = %v, want %v", got, tt.want)
				t.Error(cap(c.GetKeysBindings()))
				t.Error(cap(tt.want))

			}
		})
	}
}

func TestConfigLoader_GetKeyBind(t *testing.T) {
	type fields struct {
		Keybindings []KeyBind
		Aliases     []Alias
		Prompt      string
	}
	type args struct {
		action string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint16
	}{
		{
			name: "0",
			fields: fields{
				Keybindings: []KeyBind{{1, "asd"}},
			},
			args: args{
				action: "xcvb",
			},
			want: 0,
		},
		{
			name: "1",
			fields: fields{
				Keybindings: []KeyBind{{1, "asd"}},
			},
			args: args{
				action: "asd",
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ConfigLoader{
				Keybindings: tt.fields.Keybindings,
				Aliases:     tt.fields.Aliases,
				Prompt:      tt.fields.Prompt,
			}
			if got := c.GetKeyBind(tt.args.action); got != tt.want {
				t.Errorf("ConfigLoader.GetKeyBind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigLoader_GetEnvs(t *testing.T) {
	type fields struct {
		Keybindings []KeyBind
		Aliases     []Alias
		Prompt      string
		Envs        []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "1",
			fields: fields{
				Envs: []string{"ad"},
			},
			want: []string{"ad"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ConfigLoader{
				Keybindings: tt.fields.Keybindings,
				Aliases:     tt.fields.Aliases,
				Prompt:      tt.fields.Prompt,
				Envs:        tt.fields.Envs,
			}
			if got := c.GetEnvs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigLoader.GetEnvs() = %v, want %v", got, tt.want)
			}
		})
	}
}
