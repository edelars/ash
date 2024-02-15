package keys_bindings

import (
	"reflect"
	"testing"

	"ash/internal/commands"
	"ash/internal/dto"

	"github.com/stretchr/testify/assert"
)

func TestNewKeyBindingsManager(t *testing.T) {
	cl := confLoaderImpl{}
	cr := commRouterImpl{}

	kb := NewKeyBindingsManager(nil, &cl, cr)

	assert.Equal(t, 2, len(kb.bindings))
	assert.Equal(t, ":exec", kb.bindings[13].GetName())
	assert.Equal(t, "get", kb.bindings[3].GetName())
}

type confLoaderImpl struct{}

func (c confLoaderImpl) GetKeysBindings() []struct {
	Key    uint16
	Action string
} {
	res := []struct {
		Key    uint16
		Action string
	}{
		{
			Key:    13,
			Action: ":exec",
		},
		{
			Key:    220,
			Action: "wrong_comm",
		},
		{
			Key:    3,
			Action: "get",
		},
	}
	return res
}

type s struct {
	Key    int
	Action string
}

type commRouterImpl struct{}

func (r commRouterImpl) SearchCommands(_ dto.InternalContextIface, patterns ...dto.PatternIface) dto.CommandRouterSearchResult {
	res := commands.NewCommandRouterSearchResult()
	p := commands.NewPattern(":exec", true)
	res.AddResult(&searchResult{
		name:         ":exec",
		commandsData: []dto.CommandIface{commands.NewCommand(":exec", nil, true)},
		patternValue: p,
	})
	res.AddResult(&searchResult{
		name:         "get",
		commandsData: []dto.CommandIface{commands.NewCommand("get", nil, true)},
		patternValue: commands.NewPattern("get", true),
	})

	return res
}

type searchResult struct {
	name         string
	commandsData []dto.CommandIface
	patternValue dto.PatternIface
}

func (searchresult *searchResult) GetPriority() uint8 {
	panic("not implemented") // TODO: Implement
}

func (searchresult *searchResult) GetSourceName() string {
	return searchresult.name
}

func (searchresult *searchResult) GetCommands() []dto.CommandIface {
	return searchresult.commandsData
}

func (searchresult *searchResult) GetPattern() dto.PatternIface {
	return searchresult.patternValue
}

func (searchresult *searchResult) Founded() int {
	return len(searchresult.commandsData)
}

func TestKeyBindingsManager_GetCommandByKey(t *testing.T) {
	type fields struct {
		bindings map[uint16]dto.CommandIface
	}
	type args struct {
		key uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   dto.CommandIface
	}{
		{
			name: "ok",
			fields: fields{
				bindings: map[uint16]dto.CommandIface{13: commands.NewCommand("12", nil, true), 22: nil},
			},
			args: args{
				key: 13,
			},
			want: commands.NewCommand("12", nil, true),
		},
		{
			name: "nil",
			fields: fields{
				bindings: map[uint16]dto.CommandIface{13: commands.NewCommand("12", nil, true)},
			},
			args: args{
				key: 113,
			},
			want: nil,
		},
		{
			name: "nil 2",
			fields: fields{
				bindings: map[uint16]dto.CommandIface{},
			},
			args: args{
				key: 13,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KeyBindingsManager{
				bindings: tt.fields.bindings,
			}
			if got := k.GetCommandByKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyBindingsManager.GetCommandByKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

