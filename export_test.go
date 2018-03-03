package main

import (
	"reflect"
	"testing"

	"github.com/toolateforteddy/arbitrary/src/arbitrary"
)

func anonimize(data interface{}) interface{} {
	var anon interface{}
	err := arbitrary.Hydrate(data, &anon)
	if err != nil {
		panic(err)
	}
	return anon
}

func TestFormatForShellExport(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "easy",
			args: args{
				data: anonimize(
					map[string]map[string]string{
						"foo": map[string]string{
							"bar": "baz",
						},
					},
				),
			},
			want:    []string{"FOO_BAR=baz"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatForShellExport(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatForShellExport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FormatForShellExport() = %v, want %v", got, tt.want)
			}
		})
	}
}
