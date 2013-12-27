package xml

import (
	"reflect"
	"testing"
)

var builderBuildTests = []struct {
	builder
	want Node
}{
	{
		builder{
			name: "a",
		},
		&Element{
			Name: Name{Local: "a"},
		},
	},
}

func TestBuilderBuild(t *testing.T) {
	for i, tt := range builderBuildTests {
		got := tt.builder.build()
		if !reflect.DeepEqual(tt.want, got) {
			t.Errorf("%d: builder.build() = %#v, want %#v", i+1, got, tt.want)
		}
	}
}
