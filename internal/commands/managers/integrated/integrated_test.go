package integrated

import (
	"context"
	"reflect"
	"testing"

	"ash/internal/commands"
	"ash/internal/commands/managers/integrated/list"
	"ash/internal/internal_context"

	"github.com/stretchr/testify/assert"
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

func Test_getStepValue(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int8
	}{
		{
			name: "getexit",
			args: args{
				s: "getexit",
			},
			want: 14,
		},
		{
			name: "exit",
			args: args{
				s: "exit",
			},
			want: 25,
		},
		{
			name: "exit1",
			args: args{
				s: "exit1",
			},
			want: 20,
		},
		{
			name: "1",
			args: args{
				s: "1",
			},
			want: 100,
		},
		{
			name: "0",
			args: args{
				s: "",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStepValue(tt.args.s); got != tt.want {
				t.Errorf("getStepValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntergatedManager_searchPatternInCommands(t *testing.T) {
	f := func(ctx context.Context, internalContext internal_context.InternalContextIface, inputChan chan []byte, outputChan chan []byte) {
		panic("exit")
	}
	c1 := commands.NewCommand("exit", f)
	c2 := commands.NewCommand("dobus", f)
	c3 := commands.NewCommand("gettalk", f)

	type fields struct {
		data []commands.CommandIface
	}
	type args struct {
		searchPattern string
		founded       foundedData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   foundedData
	}{
		{
			name: "100",
			fields: fields{
				data: []commands.CommandIface{c1, c2},
			},
			args: args{
				searchPattern: "exit",
				founded:       map[commands.CommandIface]int8{},
			},
			want: map[commands.CommandIface]int8{c1: 100},
		},
		{
			name: "50",
			fields: fields{
				data: []commands.CommandIface{c1, c2},
			},
			args: args{
				searchPattern: "et",
				founded:       map[commands.CommandIface]int8{},
			},
			want: map[commands.CommandIface]int8{c1: 50},
		},
		{
			name: "none",
			fields: fields{
				data: []commands.CommandIface{c1, c2},
			},
			args: args{
				searchPattern: "555",
				founded:       map[commands.CommandIface]int8{},
			},
			want: map[commands.CommandIface]int8{},
		},
		{
			name: "44",
			fields: fields{
				data: []commands.CommandIface{c2, c3},
			},
			args: args{
				searchPattern: "ttk",
				founded:       map[commands.CommandIface]int8{},
			},
			want: map[commands.CommandIface]int8{c3: 44},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := IntergatedManager{
				data: tt.fields.data,
			}
			if got := m.searchPatternInCommands(tt.args.searchPattern, tt.args.founded); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntergatedManager.searchPatternInCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntergatedManager_SearchCommands(t *testing.T) {
	f := func(ctx context.Context, internalContext internal_context.InternalContextIface, inputChan chan []byte, outputChan chan []byte) {
		panic("exit")
	}
	c1 := commands.NewCommand("exit", f)
	c2 := commands.NewCommand("dobus", f)
	c3 := commands.NewCommand("ggalk", f)

	im := IntergatedManager{data: []commands.CommandIface{c1, c2, c3}}
	ch := make(chan commands.CommandManagerSearchResult, 1)
	p := commands.NewPattern("ext")

	im.SearchCommands(ch, p)
	res := <-ch
	assert.Equal(t, 1, len(res.GetCommands()))
	assert.Equal(t, "exit", res.GetCommands()[0].GetName())
}
