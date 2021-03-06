package combinator

import (
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

func TestCalcCLCode(t *testing.T) {
	type TD struct {
		clcode string
		cs     []Combinator
		n      int
		expect string
		desc   string
	}
	tds := []TD{
		TD{
			clcode: "Sxyz",
			cs:     cs,
			n:      -1,
			expect: "xz(yz)",
			desc:   "一度だけ計算する",
		},
		TD{
			clcode: "SKII",
			cs:     cs,
			n:      -1,
			expect: "I",
			desc:   "最後まで計算する",
		},
		TD{
			clcode: "SSSSS",
			cs:     cs,
			n:      -1,
			expect: "SS((SS)S)",
			desc:   "多段ネストの計算をする",
		},
		TD{
			clcode: "((((SSSSS))))",
			cs:     cs,
			n:      -1,
			expect: "SS((SS)S)",
			desc:   "多段ネストの計算をする",
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
			n:      -1,
			expect: "xz",
			desc:   "計算結果の関係で先頭に括弧で括られたコンビネータが来ても計算する",
		},
		TD{
			clcode: "Sxyza",
			cs:     cs,
			n:      -1,
			expect: "xz(yz)a",
			desc:   "計算結果は結合される",
		},
		TD{
			clcode: "SKIx",
			cs:     cs,
			n:      0,
			expect: "SKIx",
			desc:   "一度も計算しない",
		},
		TD{
			clcode: "SSSKS",
			cs:     cs,
			n:      1,
			expect: "SK(SK)S",
			desc:   "一度だけ計算する",
		},
	}
	for _, td := range tds {
		clcode, cs, n, expect, desc := td.clcode, td.cs, td.n, td.expect, td.desc
		actual := CalcCLCode(clcode, cs, n)
		assert.Equal(t, expect, actual, desc, clcode, cs, n)
	}
}

func TestGetBracketCombinator(t *testing.T) {
	type TD struct {
		clcode string
		expect string
	}
	tds := []TD{
		TD{clcode: "(S)KIx", expect: "(S)"},
		TD{clcode: "(SKI)x", expect: "(SKI)"},
		TD{clcode: "(SKI)", expect: "(SKI)"},
		TD{clcode: "S(SKI)", expect: "S"},
		TD{clcode: "(SKI", expect: "(SKI"},
	}
	for _, v := range tds {
		clcode, expect := v.clcode, v.expect
		got1 := getBracketCombinator(clcode)
		assert.Equal(t, expect, got1)
	}
}

func TestCalcCLCode1Time(t *testing.T) {
	assert.Equal(t, "xz(yz)", CalcCLCode1Time("Sxyz", cs))
	assert.Equal(t, "xz(yz)!", CalcCLCode1Time("Sxyz!", cs))
	assert.Equal(t, "xz(yz)", CalcCLCode1Time("(S)xyz", cs))
	assert.Equal(t, "xz(yz)xyz", CalcCLCode1Time("(Sxyz)xyz", cs))
	assert.Equal(t, "Sxy", CalcCLCode1Time("Sxy", cs))
	assert.Equal(t, "S", CalcCLCode1Time("S", cs))
	assert.Equal(t, "", CalcCLCode1Time("", cs))
	assert.Equal(t, "Sxyz", CalcCLCode1Time("Sxyz", []Combinator{}))
}

func TestTrimBracket(t *testing.T) {
	assert.Equal(t, "S", trimBracket("(S)"), "1つ括弧を外す")
	assert.Equal(t, "S", trimBracket("((((S))))"), "ネストした括弧を外す")
	assert.Equal(t, "(S)S", trimBracket("((((S)S)))"), "ネストした括弧を外す")
	assert.Equal(t, "S(S)S", trimBracket("(((S(S)S)))"), "ネストした括弧を外す")
	assert.Equal(t, "", trimBracket(""), "空のときは空を返す")
	assert.Equal(t, "(S", trimBracket("(S"), "括弧不正のときはそのまま返す")
	assert.Equal(t, "S)", trimBracket("S)"), "括弧不正のときはそのまま返す")
}

func TestCalcCombinatorArgs(t *testing.T) {
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
		actual := calcCombinatorArgs(td.inCLCode, td.inCombinator)
		assert.Equal(t, expect, actual, desc, td.inCLCode, td.inCombinator)
	}
}

func TestSplitPrefixArgsSuffixCombinators(t *testing.T) {
	var (
		pref string
		args []string
		suff string
	)

	pref, args, suff = splitPrefixArgsSuffixCombinators("Sxyz", cs)
	assert.Equal(t, "S", pref)
	assert.Equal(t, []string{"x", "y", "z"}, args)
	assert.Equal(t, "", suff)

	pref, args, suff = splitPrefixArgsSuffixCombinators("Sxyz!", cs)
	assert.Equal(t, "S", pref)
	assert.Equal(t, []string{"x", "y", "z"}, args)
	assert.Equal(t, "!", suff)

}

func TestGetPrefixCombinator(t *testing.T) {
	type TD struct {
		inCLCode      string
		inCombinators []Combinator
		expect        string
		desc          string
	}
	tds := []TD{
		TD{
			inCLCode:      "Sabc",
			inCombinators: cs,
			expect:        "S",
			desc:          "正常系",
		},
		TD{
			inCLCode:      "(abc)x",
			inCombinators: []Combinator{},
			expect:        "(abc)",
			desc:          "括弧で括られた文字列はコンビネータである",
		},
		TD{
			inCLCode:      "(ab(xzy))x",
			inCombinators: []Combinator{},
			expect:        "(ab(xzy))",
			desc:          "ネストした括弧もコンビネータである",
		},
		TD{
			inCLCode:      "(abc)(xyz)",
			inCombinators: []Combinator{},
			expect:        "(abc)",
			desc:          "括弧が連続しても別のコンビネータ",
		},
		TD{
			inCLCode:      "abc",
			inCombinators: []Combinator{Combinator{Name: "aAAAAAA"}},
			expect:        "a",
			desc:          "コンビネータにマッチしない",
		},
		TD{
			inCLCode:      "Sabc",
			inCombinators: []Combinator{Combinator{Name: "Sabc"}},
			expect:        "Sabc",
			desc:          "複数文字コンビネータ",
		},
	}
	for _, td := range tds {
		clcode, comb, desc, expect := td.inCLCode, td.inCombinators, td.desc, td.expect
		actual := getPrefixCombinator(clcode, comb)
		assert.Equal(t, expect, actual, desc, clcode, comb)
	}
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
		actual := getCombinatorArgs(clcode, cs)
		assert.Equal(t, expect, actual, desc, clcode, cs)
	}
}
