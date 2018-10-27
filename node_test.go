package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCLCode(t *testing.T) {
	type TD struct {
		in     string
		inConf Config
		expect []Node
		desc   string
	}
	tds := []TD{
		TD{
			in: "Sxyz",
			expect: []Node{
				Node{Name: "S"},
				Node{Name: "x"},
				Node{Name: "y"},
				Node{Name: "z"},
			},
			desc: "正常系",
		},
		TD{
			in:     "S",
			expect: []Node{Node{Name: "S"}},
			desc:   "一文字だけ",
		},
		TD{
			in: "(SKI)(a)(abcd)(ab(abc))",
			expect: []Node{
				Node{Name: "(SKI)", Nodes: []Node{
					Node{Name: "S"},
					Node{Name: "K"},
					Node{Name: "I"},
				}},
				Node{Name: "(a)", Nodes: []Node{
					Node{Name: "a"},
				}},
				Node{Name: "(abcd)", Nodes: []Node{
					Node{Name: "a"},
					Node{Name: "b"},
					Node{Name: "c"},
					Node{Name: "d"},
				}},
				Node{Name: "(ab(abc))", Nodes: []Node{
					Node{Name: "a"},
					Node{Name: "b"},
					Node{Name: "(abc)", Nodes: []Node{
						Node{Name: "a"},
						Node{Name: "b"},
						Node{Name: "c"},
					}},
				}},
			},
			desc: "()がネストしたCLTerm",
		},
		TD{
			in: "(SKI )(a　)(abcd\t)(ab(abc))",
			expect: []Node{
				Node{Name: "(SKI)", Nodes: []Node{
					Node{Name: "S"},
					Node{Name: "K"},
					Node{Name: "I"},
				}},
				Node{Name: "(a)", Nodes: []Node{
					Node{Name: "a"},
				}},
				Node{Name: "(abcd)", Nodes: []Node{
					Node{Name: "a"},
					Node{Name: "b"},
					Node{Name: "c"},
					Node{Name: "d"},
				}},
				Node{Name: "(ab(abc))", Nodes: []Node{
					Node{Name: "a"},
					Node{Name: "b"},
					Node{Name: "(abc)", Nodes: []Node{
						Node{Name: "a"},
						Node{Name: "b"},
						Node{Name: "c"},
					}},
				}},
			},
			desc: "空白が混在するデータ。削除されることを期待",
		},
		TD{
			in: "(SKI)",
			expect: []Node{
				Node{Name: "(SKI)", Nodes: []Node{
					Node{Name: "S"},
					Node{Name: "K"},
					Node{Name: "I"},
				}},
			},
			desc: "カッコでくくられたデータは分解されること",
		},
		TD{
			in: "Sxyz",
			inConf: Config{
				CombinatorFormat{CombinatorName: "S"},
				CombinatorFormat{CombinatorName: "K"},
				CombinatorFormat{CombinatorName: "I"},
			},
			expect: []Node{
				Node{Name: "S"},
				Node{Name: "x"},
				Node{Name: "y"},
				Node{Name: "z"},
			},
			desc: "設定ファイル指定",
		},
		TD{
			in: "Sabcz",
			inConf: Config{
				CombinatorFormat{CombinatorName: "Sabc"},
			},
			expect: []Node{
				Node{Name: "Sabc"},
				Node{Name: "z"},
			},
			desc: "複数文字のコンビネータ定義がある場合",
		},
	}
	for _, td := range tds {
		in, expect, desc, conf := td.in, td.expect, td.desc, td.inConf
		nodes, err := ParseCLCode(in, conf)
		assert.Equal(t, expect, nodes, desc, "in="+in)
		assert.NoError(t, err, "エラーが発生した", "in="+in)
	}

	// 異常系
	// エラーが返ることだけチェック
	tds = []TD{
		TD{
			in:     "SKI((A)",
			expect: nil,
			desc:   "()対応不正",
		},
		TD{
			in:     "SKI(A))",
			expect: nil,
			desc:   "()対応不正",
		},
		TD{
			in:     "SKI()",
			expect: nil,
			desc:   "()のみのデータあり",
		},
		TD{
			in:     "",
			expect: nil,
			desc:   "引数が空文字",
		},
	}
	for _, td := range tds {
		in, expect, desc := td.in, td.expect, td.desc
		nodes, err := ParseCLCode(in)
		assert.Equal(t, expect, nodes, desc, "in="+in)
		assert.Error(t, err, "エラーが発生しなかった", "in="+in)
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
			inCLCode:      "(abc)",
			inCombinators: []string{"aAAAAAA"},
			expect:        "(",
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
		actual := getPrefixCombinator(clcode, comb)
		assert.Equal(t, expect, actual, desc, fmt.Sprintf("in:{clcode:%v,comb:%v}", clcode, comb))
	}
}
