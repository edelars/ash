package data_source

import (
	"reflect"
	"testing"

	"ash/internal/dto"
)

type searchResult struct {
	name         string
	commandsData []dto.CommandIface
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
	Weight int8
	Name   string
}

func (commandimpl *commandImpl) SetMathWeight(weight int8) {
	panic("not implemented") // TODO: Implement
}

func (commandimpl *commandImpl) GetExecFunc() dto.ExecF {
	panic("not implemented") // TODO: Implement
}

func (commandimpl *commandImpl) WithArgs(args string) dto.CommandIface {
	panic("not implemented") // TODO: Implement
}

func (commandimpl *commandImpl) GetArgs() string {
	panic("not implemented") // TODO: Implement
}

func (commandImpl *commandImpl) GetMathWeight() int8 {
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
			name:   "a",
			fields: fields{},
			args: args{
				i: 0,
			},
			want: 97,
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
			want: 97,
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
					name: "1",
				}, &searchResult{
					name:         "2",
					commandsData: []dto.CommandIface{&commandImpl{Name: "22"}},
				}, &searchResult{
					name:         "3",
					commandsData: []dto.CommandIface{&commandImpl{Name: "34", Weight: 11}, &commandImpl{Name: "33", Weight: 100}},
				}}, keyMapping: make(map[rune]dto.CommandIface),
			},
			args: args{
				avalaibleSpace:         30,
				overheadLinesPerSource: 2,
			},
			want: []dto.GetDataResult{{
				SourceName: "2",
				Items: []dto.GetDataResultItem{{
					Name:       "22",
					ButtonRune: 'a',
				}},
			}, {
				SourceName: "3",
				Items: []dto.GetDataResultItem{{
					Name:       "33",
					ButtonRune: 'b',
				}, {
					Name:       "34",
					ButtonRune: 'c',
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			datasourceimpl := &dataSourceImpl{
				originalData: tt.fields.originalData,
				keyMapping:   tt.fields.keyMapping,
			}
			if got := datasourceimpl.GetData(tt.args.avalaibleSpace, tt.args.overheadLinesPerSource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceImpl.GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}
