package commands

import (
	"reflect"
	"testing"

	"ash/internal/dto"

	"github.com/stretchr/testify/assert"
)

func TestCommandRouterSearchResult_addResult(t *testing.T) {
	r := NewCommandRouterSearchResult()
	p1 := NewPattern("qweq", false)
	c := Command{weight: 44}
	s1 := searchResult{
		name:         "3333",
		commandsData: []dto.CommandIface{&c},
		patternValue: p1,
	}
	r.AddResult(&s1)

	assert.Equal(t, 1, len(r.GetDataByPattern(p1)))
	cc := r.GetDataByPattern(p1)
	assert.Equal(t, 1, len(cc))
	assert.Equal(t, uint8(44), cc[0].GetCommands()[0].GetMathWeight())

	p2 := NewPattern("qw123123", false)
	cc2 := r.GetDataByPattern(p2)
	assert.Equal(t, 0, len(cc2))
}

type CommandManagerImpl struct{}

func (commandmanagerimpl CommandManagerImpl) SearchCommands(iContext dto.InternalContextIface, resultChan chan dto.CommandManagerSearchResult, patterns ...dto.PatternIface) {
	c := Command{weight: 44}
	resultChan <- &searchResult{
		name:         "2",
		commandsData: []dto.CommandIface{&c},
		patternValue: patterns[0],
	}
}

type CommandManagerImpl2 struct{}

func (commandmanagerimpl CommandManagerImpl2) SearchCommands(iContext dto.InternalContextIface, resultChan chan dto.CommandManagerSearchResult, patterns ...dto.PatternIface) {
	resultChan <- &searchResult{
		name:         "3",
		commandsData: []dto.CommandIface{},
		patternValue: patterns[0],
	}
}

func TestCommandRouter_SearchCommands(t *testing.T) {
	cr := NewCommandRouter(CommandManagerImpl{}, CommandManagerImpl2{})
	p1 := NewPattern("123", false)

	r := cr.SearchCommands(nil, p1)
	cc := r.GetDataByPattern(p1)

	for _, v := range cc {
		assert.Equal(t, 1, v.Founded())
	}
	assert.Equal(t, 1, len(cc))
	assert.Equal(t, p1, cc[0].GetPattern())
	assert.Equal(t, "2", cc[0].GetSourceName())
	assert.Equal(t, uint8(44), cc[0].GetCommands()[0].GetMathWeight())
}

func Test_sortCommandRouterSearchResult(t *testing.T) {
	type args struct {
		cmsrs []dto.CommandManagerSearchResult
	}
	tests := []struct {
		name string
		args args
		want []dto.CommandManagerSearchResult
	}{
		{
			name: "1",
			args: args{
				cmsrs: []dto.CommandManagerSearchResult{
					&searchResult{
						name:     "1",
						priority: 4,
					},
					&searchResult{
						name:     "2",
						priority: 1,
					},
					&searchResult{
						name:     "3",
						priority: 9,
					},
					&searchResult{
						name:     "4",
						priority: 100,
					},
				},
			},
			want: []dto.CommandManagerSearchResult{
				&searchResult{
					name:     "4",
					priority: 100,
				},
				&searchResult{
					name:     "3",
					priority: 9,
				},
				&searchResult{
					name:     "1",
					priority: 4,
				},
				&searchResult{
					name:     "2",
					priority: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortCommandRouterSearchResult(tt.args.cmsrs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortCommandRouterSearchResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_commandRouterSearchResulGetDataByPattern(t *testing.T) {
	type fields struct {
		data map[dto.PatternIface][]dto.CommandManagerSearchResult
	}
	type args struct {
		pattern dto.PatternIface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []dto.CommandManagerSearchResult
	}{
		{
			name: "ok",
			fields: fields{
				data: map[dto.PatternIface][]dto.CommandManagerSearchResult{
					NewPattern("112", false): {
						&searchResult{
							name:     "1",
							priority: 4,
						},
						&searchResult{
							name:     "2",
							priority: 1,
						},
						&searchResult{
							name:     "3",
							priority: 9,
						},
						&searchResult{
							name:     "4",
							priority: 100,
						},
					},
				},
			},
			args: args{
				pattern: NewPattern("112", false),
			},
			want: []dto.CommandManagerSearchResult{
				&searchResult{
					name:     "4",
					priority: 100,
				},
				&searchResult{
					name:     "3",
					priority: 9,
				},
				&searchResult{
					name:     "1",
					priority: 4,
				},
				&searchResult{
					name:     "2",
					priority: 1,
				},
			},
		},
		{
			name: "nil",
			fields: fields{
				data: map[dto.PatternIface][]dto.CommandManagerSearchResult{
					NewPattern("112", false): {
						&searchResult{
							name:     "1",
							priority: 4,
						},
						&searchResult{
							name:     "2",
							priority: 1,
						},
						&searchResult{
							name:     "3",
							priority: 9,
						},
						&searchResult{
							name:     "4",
							priority: 100,
						},
					},
				},
			},
			args: args{
				pattern: NewPattern("zcfsd", false),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &commandRouterSearchResult{
				data: tt.fields.data,
			}
			if got := c.GetDataByPattern(tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commandRouterSearchResult.GetDataByPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
