package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cs = []Combinator{
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

func TestGetCombinatorArgs(t *testing.T) {
	type TD struct {
		clcode string
		cs     []Combinator
		expect []string
		desc   string
	}
	tds := []TD{
		TD{
			clcode: "Sxyz",
			cs:     cs,
			expect: []string{"x", "y", "z"},
			desc:   "正常系",
		},
		TD{
			clcode: "S(abc)(ab)(c)",
			cs:     cs,
			expect: []string{"(abc)", "(ab)", "(c)"},
			desc:   "括弧で括られた文字列は1コンビネータ",
		},
		TD{
			clcode: "S(abc)(ab)",
			cs:     cs,
			expect: []string{},
			desc:   "引数が不足しているときは空配列を返す",
		},
		TD{
			clcode: "",
			cs:     cs,
			expect: []string{},
			desc:   "何も渡されないときは空配列を返す",
		},
		TD{
			clcode: "Sxyz",
			cs:     []Combinator{},
			expect: []string{},
			desc:   "定義済みコンビネータがない時は何も返さない",
		},
		TD{
			clcode: "Ixyz",
			cs:     cs,
			expect: []string{"x"},
			desc:   "Iコンビネータ",
		},
		TD{
			clcode: "Zxyz",
			cs:     cs,
			expect: []string{},
			desc:   "マッチするコンビネータがない場合は空配列を返す",
		},
	}
	for _, td := range tds {
		clcode, cs, expect, desc := td.clcode, td.cs, td.expect, td.desc
		actual := GetCombinatorArgs(clcode, cs)
		assert.Equal(t, expect, actual, desc, clcode, cs)
	}
}

func TestCalcCLCode(t *testing.T) {
	type TD struct {
		clcode string
		cs     []Combinator
		expect string
		desc   string
	}
	tds := []TD{
		TD{
			clcode: "Sxyz",
			cs:     cs,
			expect: "xz(yz)",
			desc:   "一度だけ計算する",
		},
		TD{
			clcode: "SKII",
			cs:     cs,
			expect: "I",
			desc:   "最後まで計算する",
		},
		TD{
			clcode: "ZKxyz",
			cs: []Combinator{
				Combinator{
					Name:      "Z",
					ArgsCount: 1,
					Format:    "({0})",
				},
				Combinator{
					Name:      "K",
					ArgsCount: 2,
					Format:    "{0}",
				},
			},
			expect: "xz",
			desc:   "計算結果の関係で先頭に括弧で括られたコンビネータが来ても計算する",
		},
		TD{
			clcode: "Sxyza",
			cs:     cs,
			expect: "xz(yz)a",
			desc:   "計算結果は結合される",
		},
	}
	for _, td := range tds {
		clcode, cs, expect, desc := td.clcode, td.cs, td.expect, td.desc
		actual := CalcCLCode(clcode, cs)
		assert.Equal(t, expect, actual, desc, clcode, cs)
	}
}

func TestCalcHeadCombinator(t *testing.T) {
	s := Combinator{
		Name:      "S",
		ArgsCount: 3,
		Format:    "{0}{2}({1}{2})",
	}
	type TD struct {
		inCLCode     []string
		inCombinator Combinator
		out          string
		desc         string
	}
	tds := []TD{
		TD{
			inCLCode:     []string{"x", "y", "z"},
			inCombinator: s,
			out:          "xz(yz)",
			desc:         "正常系",
		},
		TD{
			inCLCode:     []string{""},
			inCombinator: s,
			out:          "",
			desc:         "空文字だけのときは空文字を返す",
		},
		TD{
			inCLCode:     []string{"x", "y"},
			inCombinator: s,
			out:          "xy",
			desc:         "引数不足は処理せず結合して返す",
		},
		TD{
			inCLCode:     []string{},
			inCombinator: s,
			out:          "",
			desc:         "引数が空のときは空文字列を返す",
		},
		TD{
			inCLCode:     nil,
			inCombinator: s,
			out:          "",
			desc:         "引数がnilのときは空文字列を返す",
		},
		TD{
			inCLCode:     []string{"x", "y"},
			inCombinator: Combinator{},
			out:          "",
			desc:         "コンビネータが空のときはそのまま返す。",
		},
	}
	for _, td := range tds {
		expect, desc := td.out, td.desc
		actual := CalcHeadCombinator(td.inCLCode, td.inCombinator)
		assert.Equal(t, expect, actual, desc, td.inCLCode, td.inCombinator)
	}
}

func TestGetPrefgixCombinator(t *testing.T) {
	type TD struct {
		inCLCode      string
		inCombinators []string
		expect        string
		desc          string
	}
	tds := []TD{
		TD{
			inCLCode:      "Sabc",
			inCombinators: []string{"S", "K", "I"},
			expect:        "S",
			desc:          "正常系",
		},
		TD{
			inCLCode:      "(abc)x",
			inCombinators: []string{},
			expect:        "(abc)",
			desc:          "括弧で括られた文字列はコンビネータである",
		},
		TD{
			inCLCode:      "(ab(xzy))x",
			inCombinators: []string{},
			expect:        "(ab(xzy))",
			desc:          "ネストした括弧もコンビネータである",
		},
		TD{
			inCLCode:      "(abc)(xyz)",
			inCombinators: []string{},
			expect:        "(abc)",
			desc:          "括弧が連続しても別のコンビネータ",
		},
		TD{
			inCLCode:      "abc",
			inCombinators: []string{"aAAAAAA"},
			expect:        "a",
			desc:          "コンビネータにマッチしない",
		},
		TD{
			inCLCode:      "Sabc",
			inCombinators: []string{"Sabc"},
			expect:        "Sabc",
			desc:          "複数文字コンビネータ",
		},
	}
	for _, td := range tds {
		clcode, comb, desc, expect := td.inCLCode, td.inCombinators, td.desc, td.expect
		actual := GetPrefixCombinator(clcode, comb)
		assert.Equal(t, expect, actual, desc, fmt.Sprintf("in:{clcode:%v,comb:%v}", clcode, comb))
	}
}
