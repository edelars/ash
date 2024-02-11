package selection_window

import (
	"reflect"
	"testing"
)

func Test_selectionWindow_calculateColumnsWidth(t *testing.T) {
	type fields struct {
		mainW                  int
		columnGap              int
		showCommandDescription bool
	}
	type args struct {
		mainFieldMaxWid int
		descFieldMaxWid int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantX    int
		wantMaxX int
	}{
		{
			name: "1",
			fields: fields{
				mainW:                  100,
				columnGap:              2,
				showCommandDescription: true,
			},
			wantX:    49,
			wantMaxX: 48,
			args: args{
				mainFieldMaxWid: 50,
				descFieldMaxWid: 60,
			},
		},
		{
			name: "2",
			fields: fields{
				mainW:                  100,
				columnGap:              2,
				showCommandDescription: true,
			},
			wantX:    0,
			wantMaxX: 0,
			args: args{
				mainFieldMaxWid: 99,
				descFieldMaxWid: 60,
			},
		},
		{
			name: "3",
			fields: fields{
				mainW:                  100,
				columnGap:              2,
				showCommandDescription: true,
			},
			wantX:    0,
			wantMaxX: 0,
			args: args{
				mainFieldMaxWid: 150,
				descFieldMaxWid: 60,
			},
		},
		{
			name: "4",
			fields: fields{
				mainW:                  100,
				columnGap:              2,
				showCommandDescription: false,
			},
			wantX:    0,
			wantMaxX: 0,
			args: args{
				mainFieldMaxWid: 50,
				descFieldMaxWid: 60,
			},
		},
		{
			name: "5",
			fields: fields{
				mainW:                  167,
				columnGap:              2,
				showCommandDescription: true,
			},
			wantX:    20,
			wantMaxX: 8,
			args: args{
				mainFieldMaxWid: 17,
				descFieldMaxWid: 8,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := &selectionWindow{
				mainW:                  tt.fields.mainW,
				columnGap:              tt.fields.columnGap,
				showCommandDescription: tt.fields.showCommandDescription,
			}
			sw.calculateColumnsWidth(tt.args.mainFieldMaxWid, tt.args.descFieldMaxWid)
			if !reflect.DeepEqual(sw.columnDescriptionX, tt.wantX) {
				t.Errorf("sw.calculateColumnsWidth columnX = %v, want %v", sw.columnDescriptionX, tt.wantX)
			}
			if !reflect.DeepEqual(sw.columnDescriptionMaxWid, tt.wantMaxX) {
				t.Errorf("sw.calculateColumnsWidth maxW= %v, want %v", sw.columnDescriptionMaxWid, tt.wantMaxX)
			}
		})
	}
}

func Test_firstN(t *testing.T) {
	type args struct {
		s string
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				s: "12345",
				n: 2,
			},
			want: "12",
		},
		{
			name: "2",
			args: args{
				s: "12345",
				n: 5,
			},
			want: "12345",
		},
		{
			name: "3",
			args: args{
				s: "",
				n: 2,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firstN(tt.args.s, tt.args.n); got != tt.want {
				t.Errorf("firstN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trimInput(t *testing.T) {
	type args struct {
		input []rune
	}
	tests := []struct {
		name string
		args args
		want []rune
	}{
		{
			name: "ls /usr",
			args: args{
				input: []rune("ls /usr"),
			},
			want: []rune("/usr"),
		},
		{
			name: "/usr",
			args: args{
				input: []rune("/usr"),
			},
			want: []rune("/usr"),
		},
		{
			name: "empty",
			args: args{
				input: []rune(""),
			},
			want: []rune(""),
		},
		{
			name: "ls -ls -a -b /usr",
			args: args{
				input: []rune("ls -ls -a -b /usr"),
			},
			want: []rune("/usr"),
		},
		{
			name: "/usr ",
			args: args{
				input: []rune("/usr "),
			},
			want: []rune(""),
		},
		{
			name: "space",
			args: args{
				input: []rune(" "),
			},
			want: []rune(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimInput(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trimInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
