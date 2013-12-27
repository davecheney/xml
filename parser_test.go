package xml

import (
	"reflect"
	"strings"
	"testing"
)

func text(s string) Node { t := Text(s); return &t }

var parserParseTests = []struct {
	input string
	want  Node
	err   error
}{
	{
		"<a></a>",
		&Element{Name: Name{Local: "a"}},
		nil,
	},
	{
		"<a/>",
		&Element{Name: Name{Local: "a"}},
		nil,
	},
	{
		"<a><b/></a>",
		&Element{Name: Name{Local: "a"},
			Children: []Node{
				&Element{Name: Name{Local: "b"}}},
		},
		nil,
	},
	{
		"<a>foo</a>",
		&Element{Name: Name{Local: "a"},
			Children: []Node{text("foo")},
		},
		nil,
	},
}

func TestParserParse(t *testing.T) {
	for i, tt := range parserParseTests {
		r := strings.NewReader(tt.input)
		p := Parser{r: r}
		got, err := p.Parse()
		if err != tt.err || !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%d: parser.Parse(%q) = %#v %v, want %#v %v", i+1, tt.input, got, err, tt.want, tt.err)
		}

	}

}
