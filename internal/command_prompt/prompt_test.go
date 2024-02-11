package command_prompt

import (
	"reflect"
	"testing"
)

func TestCommandPrompt_DeleteFromCurrentBuffer(t *testing.T) {
	type fields struct {
		currentBuffer []rune
	}
	type args struct {
		position int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    []rune
	}{
		{
			name: "1",
			fields: fields{
				currentBuffer: []rune{'1'},
			},
			args: args{
				position: 0,
			},
			wantErr: false,
			want:    []rune{},
		},
		{
			name: "err 1",
			fields: fields{
				currentBuffer: []rune{'1'},
			},
			args: args{
				position: 1,
			},
			wantErr: true,
			want:    []rune{'1'},
		},
		{
			name: "2",
			fields: fields{
				currentBuffer: []rune{'1', '2'},
			},
			args: args{
				position: 0,
			},
			wantErr: false,
			want:    []rune{'2'},
		},
		{
			name: "3",
			fields: fields{
				currentBuffer: []rune{'1', '2', '3', '4'},
			},
			args: args{
				position: 2,
			},
			wantErr: false,
			want:    []rune{'1', '2', '4'},
		},
		{
			name: "err 2",
			fields: fields{
				currentBuffer: []rune{},
			},
			args: args{
				position: 0,
			},
			wantErr: true,
			want:    []rune{},
		},
		{
			name: "4",
			fields: fields{
				currentBuffer: []rune{'1', '2', '3', '4'},
			},
			args: args{
				position: 3,
			},
			wantErr: false,
			want:    []rune{'1', '2', '3'},
		},
		{
			name: "err 3",
			fields: fields{
				currentBuffer: []rune{'1', '2', '3', '4'},
			},
			args: args{
				position: -3,
			},
			wantErr: true,
			want:    []rune{'1', '2', '3', '4'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandPrompt{
				currentBuffer: tt.fields.currentBuffer,
			}
			err := c.DeleteFromCurrentBuffer(tt.args.position)
			if err != nil != tt.wantErr {
				t.Errorf("CommandPrompt.DeleteFromCurrentBuffer() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.want, c.currentBuffer) {
				t.Errorf("CommandPrompt.DeleteFromCurrentBuffer() want = %v, got %v", tt.want, c.currentBuffer)
			}
		})
	}
}
