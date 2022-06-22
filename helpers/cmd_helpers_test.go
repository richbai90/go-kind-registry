package helpers

import (
	"reflect"
	"testing"
)

func TestCreateEnvVars(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Simple Test",
			args: args{
				args: []string {
					"a", "b", "c", "d",
				},
			},
			want: []string {
				"a=b", "c=d",
			},
			wantErr: false,
		},
		{
			
			name: "Long Test",
			args: args{
				args: []string {
					"a", "b", "c", "d", "e", "f", "g", "h",
				},
			},
			want: []string {
				"a=b", "c=d", "e=f", "g=h",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateEnvVars(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEnvVars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateEnvVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
