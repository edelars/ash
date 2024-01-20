package commands

import (
	"errors"
	"reflect"
	"testing"

	"ash/internal/dto"

	"github.com/stretchr/testify/assert"
)

func Test_commandManager_SearchCommands(t *testing.T) {
	f := func(internalC dto.InternalContextIface) {
		panic("exit")
	}
	c1 := NewCommand("exit", f)
	c2 := NewCommand("dobus", f)
	c3 := NewCommand("ggalk", f)

	im := commandManager{data: []dto.CommandIface{c1, c2, c3}}
	ch := make(chan dto.CommandManagerSearchResult, 1)
	p := NewPattern("ext", false)

	im.SearchCommands(ch, p)
	res := <-ch
	assert.Equal(t, 1, len(res.GetCommands()))
	assert.Equal(t, "exit", res.GetCommands()[0].GetName())
}

func Test_commandManager_searchPatternInCommands(t *testing.T) {
	f := func(internalC dto.InternalContextIface) {
		panic("exit")
	}
	c1 := NewCommand("exit", f)
	c2 := NewCommand("dobus", f)
	c3 := NewCommand("gettalk", f)

	type fields struct {
		data []dto.CommandIface
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
				data: []dto.CommandIface{c1, c2},
			},
			args: args{
				searchPattern: "exit",
				founded:       map[dto.CommandIface]int8{},
			},
			want: map[dto.CommandIface]int8{c1: 100},
		},
		{
			name: "50",
			fields: fields{
				data: []dto.CommandIface{c1, c2},
			},
			args: args{
				searchPattern: "et",
				founded:       map[dto.CommandIface]int8{},
			},
			want: map[dto.CommandIface]int8{c1: 50},
		},
		{
			name: "none",
			fields: fields{
				data: []dto.CommandIface{c1, c2},
			},
			args: args{
				searchPattern: "555",
				founded:       map[dto.CommandIface]int8{},
			},
			want: map[dto.CommandIface]int8{},
		},
		{
			name: "44",
			fields: fields{
				data: []dto.CommandIface{c2, c3},
			},
			args: args{
				searchPattern: "ttk",
				founded:       map[dto.CommandIface]int8{},
			},
			want: map[dto.CommandIface]int8{c3: 44},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := commandManager{
				data: tt.fields.data,
			}
			if got := m.searchPatternInCommands(tt.args.searchPattern, tt.args.founded); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntergatedManager.searchPatternInCommands() = %v, want %v", got, tt.want)
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

func Test_commandManager_precisionSearchInCommands(t *testing.T) {
	f := func(internalC dto.InternalContextIface) {
		panic("exit")
	}

	c1 := NewCommand("exit", f)
	c2 := NewCommand("dobus", f)
	c3 := NewCommand("gettalk", f)

	type fields struct {
		data []dto.CommandIface
	}
	type args struct {
		searchName string
		founded    foundedData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   foundedData
	}{
		{
			name: "exit",
			fields: fields{
				data: []dto.CommandIface{c1, c2, c3},
			},
			args: args{
				searchName: "exit",
				founded:    map[dto.CommandIface]int8{},
			},
			want: map[dto.CommandIface]int8{c1: 100},
		},
		{
			name: "gettalk",
			fields: fields{
				data: []dto.CommandIface{c1, c2, c3},
			},
			args: args{
				searchName: "gettalk",
				founded:    map[dto.CommandIface]int8{},
			},
			want: map[dto.CommandIface]int8{c3: 100},
		},
		{
			name: "none",
			fields: fields{
				data: []dto.CommandIface{c1, c2, c3},
			},
			args: args{
				searchName: "et",
				founded:    map[dto.CommandIface]int8{},
			},
			want: map[dto.CommandIface]int8{},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := commandManager{
				data: tt.fields.data,
			}
			if got := m.precisionSearchInCommands(tt.args.searchName, tt.args.founded); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntergatedManager.precisionSearchInCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchResult_GetSourceName(t *testing.T) {
	type fields struct {
		name         string
		commandsData []dto.CommandIface
		patternValue dto.PatternIface
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
				commandsData: []dto.CommandIface{},
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
		commandsData []dto.CommandIface
		patternValue dto.PatternIface
	}
	tests := []struct {
		name   string
		fields fields
		want   dto.PatternIface
	}{
		{
			name: "",
			fields: fields{
				name:         "",
				commandsData: []dto.CommandIface{},
				patternValue: NewPattern("333", false),
			},
			want: NewPattern("333", false),
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

func NewExitCommand() *Command {
	return NewCommand("exit",
		func(internalC dto.InternalContextIface) {
			internalC.GetErrChan() <- errors.New("ash exiting")
		})
}

func Test_searchResult_GetCommands(t *testing.T) {
	type fields struct {
		name         string
		commandsData []dto.CommandIface
		patternValue dto.PatternIface
	}
	c := NewExitCommand()
	tests := []struct {
		name   string
		fields fields
		want   []dto.CommandIface
	}{
		{
			name: "1",
			fields: fields{
				name:         "",
				commandsData: []dto.CommandIface{c},
				patternValue: nil,
			},
			want: []dto.CommandIface{c},
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
