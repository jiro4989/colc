package combinator

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var cs = Combinators{
	Combinator{
		Name:      "S",
		ArgsCount: 3,
		Format:    "{0}{2}({1}{2})",
	},
	Combinator{
		Name:      "K",
		ArgsCount: 2,
		Format:    "{0}",
	},
	Combinator{
		Name:      "I",
		ArgsCount: 1,
		Format:    "{0}",
	},
}

func TestParse(t *testing.T) {
	type TD struct {
		code   string
		expect Node
	}

	tds := []TD{
		TD{
			code: "Sxyz",
			expect: Node{
				Text: "Sxyz",
				Nodes: []*Node{
					&Node{Text: "S"},
					&Node{Text: "x"},
					&Node{Text: "y"},
					&Node{Text: "z"},
				},
			},
		},
		TD{
			code: "S(SS)(KI)z",
			expect: Node{
				Text: "S(SS)(KI)z",
				Nodes: []*Node{
					&Node{Text: "S"},
					&Node{Text: "SS", Nodes: []*Node{
						&Node{Text: "S"},
						&Node{Text: "S"},
					}},
					&Node{Text: "KI", Nodes: []*Node{
						&Node{Text: "K"},
						&Node{Text: "I"},
					}},
					&Node{Text: "z"},
				},
			},
		},
	}
	for _, v := range tds {
		code, expect := v.code, v.expect
		got := Parse(code, cs)
		if diff := cmp.Diff(expect, got); diff != "" {
			t.Errorf("node1 differ: (-expect +got)\n%s", diff)
		}
	}
}
