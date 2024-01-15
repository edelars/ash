package integrated

import (
	"reflect"
	"testing"

	"ash/internal/commands"
	"ash/internal/commands/managers/integrated/list"
)

func Test_searchResult_GetSourceName(t *testing.T) {
	type fields struct {
		name         string
		commandsData []commands.CommandIface
		patternValue commands.PatternIface
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "qwe",
			fields: fields{
				name:         "qwe",
				commandsData: []commands.CommandIface{},
				patternValue: nil,
			},
			want: "qwe",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchresult := &searchResult{
				name:         tt.fields.name,
				commandsData: tt.fields.commandsData,
				patternValue: tt.fields.patternValue,
			}
			if got := searchresult.GetSourceName(); got != tt.want {
				t.Errorf("searchResult.GetSourceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchResult_GetPattern(t *testing.T) {
	type fields struct {
		name         string
		commandsData []commands.CommandIface
		patternValue commands.PatternIface
	}
	tests := []struct {
		name   string
		fields fields
		want   commands.PatternIface
	}{
		{
			name: "",
			fields: fields{
				name:         "",
				commandsData: []commands.CommandIface{},
				patternValue: commands.NewPattern("333"),
			},
			want: commands.NewPattern("333"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchresult := &searchResult{
				name:         tt.fields.name,
				commandsData: tt.fields.commandsData,
				patternValue: tt.fields.patternValue,
			}
			if got := searchresult.GetPattern(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchResult.GetPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchResult_GetCommands(t *testing.T) {
	type fields struct {
		name         string
		commandsData []commands.CommandIface
		patternValue commands.PatternIface
	}
	c := list.NewExitCommand()
	tests := []struct {
		name   string
		fields fields
		want   []commands.CommandIface
	}{
		{
			name: "1",
			fields: fields{
				name:         "",
				commandsData: []commands.CommandIface{c},
				patternValue: nil,
			},
			want: []commands.CommandIface{c},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchresult := &searchResult{
				name:         tt.fields.name,
				commandsData: tt.fields.commandsData,
				patternValue: tt.fields.patternValue,
			}
			if got := searchresult.GetCommands(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchResult.GetCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}
