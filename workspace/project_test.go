package workspace

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestLoadDotFile(t *testing.T) {
	tests := []struct {
		name    string
		rd      io.Reader
		want    *DotFile
		wantErr bool
	}{
		{
			"simple",
			strings.NewReader(`name: "testproject"`),
			&DotFile{
				Name: "testproject",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadDotFile(tt.rd)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadDotFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadDotFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
