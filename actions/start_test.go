package actions

import (
	"testing"

	"github.com/richbai90/bundle-containers/helpers"
	"github.com/spf13/cobra"
)

func TestStart(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "NoArgs",
			args: args{
				cmd: helpers.NewTestCommand(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Start(tt.args.cmd, tt.args.args)
		})
	}
}

func TestGetKindVersion(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "latest version",
			want: "1.24.0@sha256:0866296e693efe1fed79d5e6c7af8df71fc73ae45e3679af05342239cdc5bc8e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetKind(""); got != tt.want {
				t.Errorf("GetKindVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
