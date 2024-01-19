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
				Keybindings: []KeyBind{{13, ":Execute"}, {14, ":Autocomplete"}},
				Aliases:     []Alias{},
				Prompt:      "ASH> ",
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
