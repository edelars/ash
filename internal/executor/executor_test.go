package executor

import (
	"context"
	"reflect"
	"testing"

	"ash/internal/commands"
	"ash/internal/dto"
	"ash/internal/internal_context"

	"github.com/stretchr/testify/assert"
)

func Test_splitToArray(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name     string
		args     args
		wantRes  []dto.PatternIface
		wantArgs []string
	}{
		{
			name: "git clone http://ya.ru | grep ok | >> out.txt",
			args: args{
				s: "git clone http://ya.ru | grep ok | >> out.txt",
			},
			wantRes:  []dto.PatternIface{commands.NewPattern("git", true), commands.NewPattern("grep", true), commands.NewPattern(">>", true)},
			wantArgs: []string{"clone http://ya.ru", "ok", "out.txt"},
		},
		{
			name: "git clone http://ya.ru|grep ok| >> out.txt",
			args: args{
				s: "git clone http://ya.ru|grep ok| >> out.txt",
			},
			wantRes:  []dto.PatternIface{commands.NewPattern("git", true), commands.NewPattern("grep", true), commands.NewPattern(">>", true)},
			wantArgs: []string{"clone http://ya.ru", "ok", "out.txt"},
		},

		{
			name: "one",
			args: args{
				s: "one",
			},
			wantRes:  []dto.PatternIface{commands.NewPattern("one", true)},
			wantArgs: []string{""},
		},
		{
			name: "one two",
			args: args{
				s: "one two",
			},
			wantRes:  []dto.PatternIface{commands.NewPattern("one", true)},
			wantArgs: []string{"two"},
		},
		{
			name: "one two |",
			args: args{
				s: "one two |",
			},
			wantRes:  []dto.PatternIface{commands.NewPattern("one", true)},
			wantArgs: []string{"two"},
		},
		{
			name: "one two | ",
			args: args{
				s: "one two | ",
			},
			wantRes:  []dto.PatternIface{commands.NewPattern("one", true)},
			wantArgs: []string{"two"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, gotArgs := splitToArrays(tt.args.s)
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("splitToArray() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("splitToArray() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestCommandExecutor_prepareExecutionList(t *testing.T) {
	cr := commRouterImpl{}
	type fields struct {
		commandRouter     routerIface
		keyBindingManager keyBindingsIface
	}
	type args struct {
		internalC dto.InternalContextIface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.InternalContextIface
		wantErr bool
	}{
		{
			name: "get | put",
			fields: fields{
				commandRouter:     cr,
				keyBindingManager: nil,
			},
			args: args{
				internalC: internal_context.InternalContext{}.WithCurrentInputBuffer([]byte("get | put")),
			},
			want:    internal_context.InternalContext{}.WithExecutionList([]dto.CommandIface{commands.NewCommand("get", nil), commands.NewCommand("put", nil)}),
			wantErr: false,
		},
		{
			name: "get |",
			fields: fields{
				commandRouter:     cr,
				keyBindingManager: nil,
			},
			args: args{
				internalC: internal_context.InternalContext{}.WithCurrentInputBuffer([]byte("get |")),
			},
			want:    internal_context.InternalContext{}.WithExecutionList([]dto.CommandIface{commands.NewCommand("get", nil)}),
			wantErr: false,
		},
		{
			name: "get asd | put 456",
			fields: fields{
				commandRouter:     cr,
				keyBindingManager: nil,
			},
			args: args{
				internalC: internal_context.InternalContext{}.WithCurrentInputBuffer([]byte("get asd | put 456")),
			},
			want:    internal_context.InternalContext{}.WithExecutionList([]dto.CommandIface{commands.NewCommand("get", nil).WithArgs("asd"), commands.NewCommand("put", nil).WithArgs("456")}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := commandExecutor{
				commandRouter:     tt.fields.commandRouter,
				keyBindingManager: tt.fields.keyBindingManager,
			}
			got, err := r.prepareExecutionList(tt.args.internalC)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandExecutor.prepareExecutionList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.GetExecutionList(), tt.want.GetExecutionList()) {
				t.Errorf("CommandExecutor.prepareExecutionList() = %v, want %v", got.GetExecutionList(), tt.want.GetExecutionList())
			}
		})
	}
}

type commRouterImpl struct{}

func (r commRouterImpl) SearchCommands(patterns ...dto.PatternIface) dto.CommandRouterSearchResult {
	res := commands.NewCommandRouterSearchResult()

	res.AddResult(&searchResult{
		name:         "get",
		commandsData: []dto.CommandIface{commands.NewCommand("get", nil)},
		patternValue: commands.NewPattern("get", true),
	})
	res.AddResult(&searchResult{
		name:         "put",
		commandsData: []dto.CommandIface{commands.NewCommand("put", nil)},
		patternValue: commands.NewPattern("put", true),
	})

	return &res
}

type searchResult struct {
	name         string
	commandsData []dto.CommandIface
	patternValue dto.PatternIface
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

func TestCommandExecutor_Execute(t *testing.T) {
	cr := commRouterImpl{}
	kb := keyBinderImpl{}
	ce := NewCommandExecutor(cr, &kb)
	ic := internal_context.NewInternalContext(context.Background(), nil, nil, nil).WithLastKeyPressed(byte(13)).WithCurrentInputBuffer([]byte("get"))
	ce.Execute(ic)
	assert.Equal(t, true, kb.Success)
}

type keyBinderImpl struct {
	Success bool
}

func (kb *keyBinderImpl) GetCommandByKey(key int) dto.CommandIface {
	if key == 13 {
		return commands.NewCommand("get", func(_ dto.InternalContextIface) {
			kb.Success = true
		})
	}
	return nil
}
