package envs_loader

import "testing"

func TestParseEnvString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantA   string
		wantB   string
		wantErr bool
	}{
		{
			name: "ASD=123",
			args: args{
				s: "ASD=123",
			},
			wantA:   "ASD",
			wantB:   "123",
			wantErr: false,
		},
		{
			name: "ASD = 123",
			args: args{
				s: "ASD = 123",
			},
			wantA:   "ASD",
			wantB:   "123",
			wantErr: false,
		},
		{
			name: "ASD123",
			args: args{
				s: "ASD123",
			},
			wantA:   "",
			wantB:   "",
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				s: "",
			},
			wantA:   "",
			wantB:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotB, err := ParseEnvString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotA != tt.wantA {
				t.Errorf("parse() gotA = %v, want %v", gotA, tt.wantA)
			}
			if gotB != tt.wantB {
				t.Errorf("parse() gotB = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
