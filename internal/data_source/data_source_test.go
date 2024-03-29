package data_source

import (
	"reflect"
	"testing"

	"ash/internal/dto"

	"github.com/stretchr/testify/assert"
)

type searchResult struct {
	name         string
	commandsData []dto.CommandIface
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

func (searchresult *searchResult) Founded() int {
	return len(searchresult.commandsData)
}

func (searchresult *searchResult) GetPattern() dto.PatternIface {
	return nil
}

type commandImpl struct {
	Weight      uint8
	Name        string
	DispalyName string
	Description string
}

func (commandimpl *commandImpl) GetDescription() string {
	return commandimpl.Description
}

func (commandimpl *commandImpl) GetDisplayName() string {
	if commandimpl.DispalyName == "" {
		return commandimpl.Name
	}
	return commandimpl.DispalyName
}

func (commandimpl *commandImpl) SetDisplayName(displayName string) {
	panic("not implemented") // TODO: Implement
}

func (commandimpl *commandImpl) SetMathWeight(weight uint8) {
	panic("not implemented") // TODO: Implement
}

func (commandimpl *commandImpl) GetExecFunc() dto.ExecutionFunction {
	panic("not implemented") // TODO: Implement
}

func (commandimpl *commandImpl) WithArgs(args []string) dto.CommandIface {
	panic("not implemented") // TODO: Implement
}

func (commandimpl *commandImpl) GetArgs() []string {
	panic("not implemented") // TODO: Implement
}

func (commandImpl *commandImpl) GetMathWeight() uint8 {
	return commandImpl.Weight
}

func (commandImpl *commandImpl) GetName() string {
	return commandImpl.Name
}

func (commandImpl *commandImpl) MustPrepareExecutionList() bool {
	return true
}

func Test_sortSlice(t *testing.T) {
	c1 := commandImpl{Weight: 99, Name: "99"}
	c2 := commandImpl{Weight: 30, Name: "30"}
	c3 := commandImpl{Weight: 10, Name: "10"}
	type args struct {
		cmds []dto.CommandIface
	}
	tests := []struct {
		name string
		args args
		want []dto.CommandIface
	}{
		{
			name: "1",
			args: args{
				cmds: []dto.CommandIface{&c1, &c3, &c2},
			},
			want: []dto.CommandIface{&c1, &c2, &c3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortSlice(tt.args.cmds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceImpl_initGetDataResult(t *testing.T) {
	type fields struct {
		originalData []dto.CommandManagerSearchResult
	}
	type args struct {
		avalaibleSpace         int
		overheadLinesPerSource int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []dto.GetDataResult
	}{
		{
			name: "1",
			fields: fields{
				originalData: []dto.CommandManagerSearchResult{&searchResult{
					name:         "1",
					commandsData: []dto.CommandIface{&commandImpl{}, &commandImpl{}, &commandImpl{}, &commandImpl{}},
				}, &searchResult{
					name:         "2",
					commandsData: []dto.CommandIface{&commandImpl{}},
				}, &searchResult{
					name:         "3",
					commandsData: []dto.CommandIface{&commandImpl{}},
				}},
			},
			args: args{
				avalaibleSpace:         20,
				overheadLinesPerSource: 2,
			},
			want: []dto.GetDataResult{{
				SourceName: "1",
				Items:      make([]dto.GetDataResultItem, 4),
			}, {
				SourceName: "2",
				Items:      make([]dto.GetDataResultItem, 1),
			}, {
				SourceName: "3",
				Items:      make([]dto.GetDataResultItem, 1),
			}},
		},
		{
			name: "2",
			fields: fields{
				originalData: []dto.CommandManagerSearchResult{&searchResult{
					name:         "1",
					commandsData: []dto.CommandIface{&commandImpl{}, &commandImpl{}, &commandImpl{}, &commandImpl{}},
				}, &searchResult{
					name:         "2",
					commandsData: []dto.CommandIface{&commandImpl{}},
				}, &searchResult{
					name:         "3",
					commandsData: []dto.CommandIface{&commandImpl{}},
				}},
			},
			args: args{
				avalaibleSpace:         9,
				overheadLinesPerSource: 2,
			},
			want: []dto.GetDataResult{{
				SourceName: "1",
				Items:      make([]dto.GetDataResultItem, 1),
			}, {
				SourceName: "2",
				Items:      make([]dto.GetDataResultItem, 1),
			}, {
				SourceName: "3",
				Items:      make([]dto.GetDataResultItem, 1),
			}},
		},
		{
			name: "3",
			fields: fields{
				originalData: []dto.CommandManagerSearchResult{&searchResult{
					name:         "1",
					commandsData: []dto.CommandIface{&commandImpl{}, &commandImpl{}, &commandImpl{}, &commandImpl{}},
				}, &searchResult{
					name:         "2",
					commandsData: []dto.CommandIface{&commandImpl{}, &commandImpl{}, &commandImpl{}, &commandImpl{}},
				}, &searchResult{
					name:         "3",
					commandsData: []dto.CommandIface{&commandImpl{}},
				}},
			},
			args: args{
				avalaibleSpace:         13,
				overheadLinesPerSource: 2,
			},
			want: []dto.GetDataResult{{
				SourceName: "1",
				Items:      make([]dto.GetDataResultItem, 2),
			}, {
				SourceName: "2",
				Items:      make([]dto.GetDataResultItem, 2),
			}, {
				SourceName: "3",
				Items:      make([]dto.GetDataResultItem, 1),
			}},
		},
		{
			name: "4",
			fields: fields{
				originalData: []dto.CommandManagerSearchResult{&searchResult{
					name:         "1",
					commandsData: []dto.CommandIface{&commandImpl{}, &commandImpl{}, &commandImpl{}, &commandImpl{}, &commandImpl{}, &commandImpl{}, &commandImpl{}, &commandImpl{}},
				}, &searchResult{
					name:         "2",
					commandsData: []dto.CommandIface{&commandImpl{}, &commandImpl{}, &commandImpl{}, &commandImpl{}},
				}, &searchResult{
					name:         "3",
					commandsData: []dto.CommandIface{&commandImpl{}, &commandImpl{}, &commandImpl{}, &commandImpl{}, &commandImpl{}},
				}, &searchResult{}, &searchResult{}, &searchResult{
					name:         "4",
					commandsData: []dto.CommandIface{&commandImpl{}},
				}},
			},
			args: args{
				avalaibleSpace:         21,
				overheadLinesPerSource: 2,
			},
			want: []dto.GetDataResult{{
				SourceName: "1",
				Items:      make([]dto.GetDataResultItem, 3),
			}, {
				SourceName: "2",
				Items:      make([]dto.GetDataResultItem, 3),
			}, {
				SourceName: "3",
				Items:      make([]dto.GetDataResultItem, 3),
			}, {
				SourceName: "4",
				Items:      make([]dto.GetDataResultItem, 1),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			datasourceimpl := &dataSourceImpl{
				originalData: tt.fields.originalData,
			}
			if got := datasourceimpl.initGetDataResult(tt.args.avalaibleSpace, tt.args.overheadLinesPerSource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceImpl.initdto.GetDataResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceImpl_generateRune(t *testing.T) {
	type fields struct {
		originalData []dto.CommandManagerSearchResult
	}
	type args struct {
		i rune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   rune
	}{
		{
			name:   "0",
			fields: fields{},
			args: args{
				i: 0,
			},
			want: 48,
		},
		{
			name:   "98",
			fields: fields{},
			args: args{
				i: 97,
			},
			want: 98,
		},
		{
			name:   "333",
			fields: fields{},
			args: args{
				i: 333,
			},
			want: 0,
		},
		{
			name:   "57",
			fields: fields{},
			args: args{
				i: 57,
			},
			want: 97,
		},
		{
			name:   "122",
			fields: fields{},
			args: args{
				i: 122,
			},
			want: 65,
		},
		{
			name:   "67",
			fields: fields{},
			args: args{
				i: 67,
			},
			want: 68,
		},
		{
			name:   "A",
			fields: fields{},
			args: args{
				i: 65,
			},
			want: 66,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			datasourceimpl := &dataSourceImpl{
				originalData: tt.fields.originalData,
			}
			if got := datasourceimpl.generateRune(tt.args.i); got != tt.want {
				t.Errorf("dataSourceImpl.generateRune() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceImpl_GetData(t *testing.T) {
	type fields struct {
		originalData []dto.CommandManagerSearchResult
		keyMapping   map[rune]dto.CommandIface
	}
	type args struct {
		avalaibleSpace         int
		overheadLinesPerSource int
	}
	type w struct {
		res []dto.GetDataResult
		m   int
		d   int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   w
	}{
		{
			name: "1",
			fields: fields{
				originalData: []dto.CommandManagerSearchResult{&searchResult{
					name: "1",
				}, &searchResult{
					name:         "2",
					commandsData: []dto.CommandIface{&commandImpl{Name: "22", Description: "123456"}},
				}, &searchResult{
					name:         "3",
					commandsData: []dto.CommandIface{&commandImpl{Name: "34", Weight: 11, Description: "4"}, &commandImpl{Name: "33", Weight: 100}, &commandImpl{Name: "3533", Weight: 90, Description: "444"}},
				}}, keyMapping: make(map[rune]dto.CommandIface),
			},
			args: args{
				avalaibleSpace:         30,
				overheadLinesPerSource: 2,
			},
			want: w{[]dto.GetDataResult{{
				SourceName: "2",
				Items: []dto.GetDataResultItem{{
					Name:        "22",
					DisplayName: "22",
					ButtonRune:  '0',
					Description: "123456",
				}},
			}, {
				SourceName: "3",
				Items: []dto.GetDataResultItem{{
					Name:        "33",
					DisplayName: "33",
					ButtonRune:  '1',
				}, {
					Name:        "3533",
					DisplayName: "3533",
					ButtonRune:  '2',
					Description: "444",
				}, {
					Name:        "34",
					DisplayName: "34",
					ButtonRune:  '3',
					Description: "4",
				}},
			}}, 4, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			datasourceimpl := &dataSourceImpl{
				originalData: tt.fields.originalData,
				keyMapping:   tt.fields.keyMapping,
			}
			got, m, d := datasourceimpl.GetData(tt.args.avalaibleSpace, tt.args.overheadLinesPerSource)
			if !reflect.DeepEqual(got, tt.want.res) {
				t.Errorf("dataSourceImpl.GetData() = %v, want %v", got, tt.want.res)
			}
			if !reflect.DeepEqual(m, tt.want.m) {
				t.Errorf("dataSourceImpl.GetData() = %v, want %v", m, tt.want.m)
			}
			if !reflect.DeepEqual(d, tt.want.d) {
				t.Errorf("dataSourceImpl.GetData() = %v, want %v", d, tt.want.d)
			}
		})
	}
}

func Test_dataSourceImpl_GetCommand(t *testing.T) {
	h := dataSourceImpl{keyMapping: make(map[rune]dto.CommandIface)}
	h.originalData = []dto.CommandManagerSearchResult{&searchResult{
		name: "1",
	}, &searchResult{
		name:         "2",
		commandsData: []dto.CommandIface{&commandImpl{Name: "22"}},
	}, &searchResult{
		name:         "3",
		commandsData: []dto.CommandIface{&commandImpl{Name: "34", Weight: 11}, &commandImpl{Name: "33", Weight: 100}},
	}}

	res, _, _ := h.GetData(12, 1)
	for r := 0; r < len(res); r++ {
		for i := 0; i < len(res[r].Items); i++ {
			assert.Equal(t, res[r].Items[i].Name, h.GetCommand(res[r].Items[i].ButtonRune).GetName())
		}
	}
}
