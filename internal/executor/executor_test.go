package executor

import (
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
	cr2 := commRouterImpl2{}
	// cr3 := commRouterImpl3{}
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
				internalC: internal_context.InternalContext{}.WithCurrentInputBuffer([]rune("get | put")),
			},
			want:    internal_context.InternalContext{}.WithExecutionList([]dto.CommandIface{commands.NewCommand("get", nil, true), commands.NewCommand("put", nil, true)}),
			wantErr: false,
		},
		{
			name: "get |",
			fields: fields{
				commandRouter:     cr,
				keyBindingManager: nil,
			},
			args: args{
				internalC: internal_context.InternalContext{}.WithCurrentInputBuffer([]rune("get |")),
			},
			want:    internal_context.InternalContext{}.WithExecutionList([]dto.CommandIface{commands.NewCommand("get", nil, true)}),
			wantErr: false,
		},
		{
			name: "get asd | put 456",
			fields: fields{
				commandRouter:     cr,
				keyBindingManager: nil,
			},
			args: args{
				internalC: internal_context.InternalContext{}.WithCurrentInputBuffer([]rune("get asd | put 456")),
			},
			want:    internal_context.InternalContext{}.WithExecutionList([]dto.CommandIface{commands.NewCommand("get", nil, true).WithArgs([]string{"asd"}), commands.NewCommand("put", nil, true).WithArgs([]string{"456"})}),
			wantErr: false,
		},
		{
			name: "error 0",
			fields: fields{
				commandRouter: cr2,
			},
			args: args{
				internalC: internal_context.InternalContext{}.WithCurrentInputBuffer([]rune("get asd | put 456")),
			},
			want:    internal_context.InternalContext{},
			wantErr: true,
		},
		// {
		// 	name: "error 2",
		// 	fields: fields{
		// 		commandRouter: cr3,
		// 	},
		// 	args: args{
		// 		internalC: internal_context.InternalContext{}.WithCurrentInputBuffer([]rune("get asd")),
		// 	},
		// 	want:    internal_context.InternalContext{},
		// 	wantErr: true,
		// },
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

func (r commRouterImpl) SearchCommands(_ dto.InternalContextIface, patterns ...dto.PatternIface) dto.CommandRouterSearchResult {
	res := commands.NewCommandRouterSearchResult()

	res.AddResult(&searchResult{
		name:         "get",
		commandsData: []dto.CommandIface{commands.NewCommand("get", nil, true)},
		patternValue: commands.NewPattern("get", true),
	})
	res.AddResult(&searchResult{
		name:         "put",
		commandsData: []dto.CommandIface{commands.NewCommand("put", nil, true)},
		patternValue: commands.NewPattern("put", true),
	})
	return res
}

type commRouterImpl2 struct{}

func (r commRouterImpl2) SearchCommands(_ dto.InternalContextIface, patterns ...dto.PatternIface) dto.CommandRouterSearchResult {
	res := commands.NewCommandRouterSearchResult()
	return res
}

type commRouterImpl3 struct{}

func (r commRouterImpl3) SearchCommands(_ dto.InternalContextIface, patterns ...dto.PatternIface) dto.CommandRouterSearchResult {
	res := commands.NewCommandRouterSearchResult()

	res.AddResult(&searchResult{
		name:         "get",
		commandsData: []dto.CommandIface{commands.NewCommand("get", nil, true)},
		patternValue: commands.NewPattern("get", true),
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
	priority     uint8
}

func (searchresult *searchResult) GetPriority() uint8 {
	return searchresult.priority
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

func TestCommandExecutor_Execute(t *testing.T) {
	cr := commRouterImpl{}
	kb := keyBinderImpl{}
	ce := NewCommandExecutor(cr, &kb)

	ic := internal_context.NewInternalContext(nil, nil, func(msg string) {}, nil, nil, nil, nil).WithLastKeyPressed(13).WithCurrentInputBuffer([]rune("get"))

	res := ce.Execute(ic)
	assert.Equal(t, true, kb.Success)
	assert.Equal(t, 0, len(ic.GetExecutionList()))

	assert.Equal(t, dto.CommandExecResultMainExit, res)
}

type keyBinderImpl struct {
	Success bool
}

func (kb *keyBinderImpl) GetCommandByKey(key uint16) dto.CommandIface {
	if key == 13 {
		return commands.NewCommand("get", func(_ dto.InternalContextIface, _ []string) dto.ExecResult {
			kb.Success = true
			return dto.CommandExecResultMainExit
		}, true)
	}
	return nil
}

func Test_splitArgsStringToArr(t *testing.T) {
	type args struct {
		a string
	}
	tests := []struct {
		name    string
		args    args
		wantRes []string
	}{
		{
			name: "1",
			args: args{
				a: " -l  -a ",
			},
			wantRes: []string{"-l", "-a"},
		},
		{
			name: "2",
			args: args{
				a: "",
			},
			wantRes: nil,
		},
		{
			name: "3",
			args: args{
				a: "la-aa",
			},
			wantRes: []string{"la-aa"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := splitArgsStringToArr(tt.args.a); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("splitArgsStringToArr() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
