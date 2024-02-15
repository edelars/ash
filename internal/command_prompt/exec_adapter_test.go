package command_prompt

import (
	"testing"

	"ash/internal/dto"
	"ash/internal/internal_context"

	"github.com/stretchr/testify/assert"
)

func Test_extractCmd(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				s: "%(cmd)",
			},
			want:    "cmd",
			wantErr: false,
		},
		{
			name: "err",
			args: args{
				s: "cmd",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractCmd(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_execAdapter_ExecCmd(t *testing.T) {
	e := execImpl2{}
	ictx := internal_context.InternalContext{}
	h := execAdapter{exec: &e, key: 13}

	// 1 test
	res, err := h.ExecCmd(ictx, "cmd")

	assert.Error(t, err)
	assert.Equal(t, "", res)

	// 2 test
	res, err = h.ExecCmd(ictx, "%(cmd)")

	assert.NoError(t, err)
	assert.Equal(t, "qqq", res)

	// 3 test
	h.key = 11
	res, err = h.ExecCmd(ictx, "%(cmd)")

	assert.Error(t, err)
	assert.Equal(t, "", res)
}

type execImpl2 struct{}

func (execimpl2 *execImpl2) Execute(internalC dto.InternalContextIface) dto.ExecResult {
	if internalC.GetLastKeyPressed() == byte(uint16(13)) {
		internalC.GetOutputWriter().Write([]byte("q\nqq\n"))
		return dto.CommandExecResultStatusOk
	} else {
		return dto.CommandExecResultNotDoAnyting
	}
}
