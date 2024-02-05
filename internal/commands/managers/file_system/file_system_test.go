package file_system

import (
	"reflect"
	"testing"
)

func Test_preparePathArr(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "1",
			args: args{
				path: "/opt/homebrew/bin:/opt/homebrew/sbin:/usr/local/bin:/System/Cryptexes/App/usr/bin:/usr/bin:/bin:/usr/sbin:",
			},
			want: []string{"/opt/homebrew/bin", "/opt/homebrew/sbin", "/usr/local/bin", "/System/Cryptexes/App/usr/bin", "/usr/bin", "/bin", "/usr/sbin"},
		},
		{
			name: "2",
			args: args{
				path: "/opt/homebrew/bin:/opt/homebrew/sbin:/usr/local/bin:/System/Cryptexes/App/usr/bin:/usr/bin:/bin:/usr/sbin",
			},
			want: []string{"/opt/homebrew/bin", "/opt/homebrew/sbin", "/usr/local/bin", "/System/Cryptexes/App/usr/bin", "/usr/bin", "/bin", "/usr/sbin"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := preparePathArr(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("preparePathArr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFileNamesInDirs(t *testing.T) {
	type args struct {
		dirs       []string
		skipDirs   bool
		searchFunc func(dir string, skipDirs bool) []string
	}
	tests := []struct {
		name    string
		args    args
		wantRes []filesResult
	}{
		{
			name: "",
			args: args{
				dirs:     []string{"1"},
				skipDirs: false,
				searchFunc: func(dir string, skipDirs bool) []string {
					return []string{"a1", "a2"}
				},
			},
			wantRes: []filesResult{{
				dir:   "1",
				files: []string{"a1", "a2"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := getFileNamesInDirs(tt.args.dirs, tt.args.skipDirs, tt.args.searchFunc); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("getFileNamesInDirs() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
