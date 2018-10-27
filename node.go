package main

import (
	"fmt"
	"strings"
)

// Node はコンビネータの名前と、その子ノードを持つ。
type Node struct {
	Name  string
	Nodes []Node
}

type Combinator struct {
	Name      string
	ArgsCount int
	Format    string
}

func CalcCLCode(clcode string, cs []Combinator) string {
	// コンビネータリストから名前だけのコンビネータを生成
	var cns []string
	for _, c := range cs {
		cns = append(cns, c.Name)
	}

	bef, aft := clcode, clcode
	for {
		if bef != aft {
			break
		}
		bef = aft
		pref := getPrefixCombinator(bef, cns)

		// 先頭コンビネータが定義済みコンビネータの中にあればセット
		var co Combinator
		for _, c := range cs {
			if c.Name == pref {
				co = c
				break
			}
		}

		t := bef[len(pref):]
		var args []string
		for i := 0; i < co.ArgsCount; i++ {
			c := getPrefixCombinator(t, cns)
			args = append(args, c)
			t = t[len(c):]
		}

		res := CalcHeadCombinator(args, co)
		aft = res
	}
	return aft
}

func CalcHeadCombinator(cs []string, co Combinator) string {
	max := co.ArgsCount
	if len(cs) < max {
		return strings.Join(cs, "")
	}

	s := co.Format
	for i := 0; i < max; i++ {
		f := fmt.Sprintf("{%d}", i)
		s = strings.Replace(s, f, cs[i], -1)
	}

	return s
}

func GetCombinatorArgs(clcode string, cs []Combinator) []string {
	var cns []string
	for _, c := range cs {
		cns = append(cns, c.Name)
	}

	pref := getPrefixCombinator(clcode, cns)

	// 先頭コンビネータが定義済みコンビネータの中にあればセット
	var co Combinator
	for _, c := range cs {
		if c.Name == pref {
			co = c
			break
		}
	}

	clcode = clcode[len(co.Name):]
	var args []string
	for i := 0; i < co.ArgsCount; i++ {
		c := getPrefixCombinator(clcode, cns)
		args = append(args, c)
		clcode = clcode[len(c):]
	}

	return args
}

// getPrefixCombinator はCLCodeの先頭のコンビネータを返す。
// 引数に渡している定義済みコンビネータが存在した場合、複数文字でも返す。
func getPrefixCombinator(clcode string, cs []string) string {
	if len(clcode) < 1 {
		return ""
	}

	var (
		ret   string
		pref  = clcode[:1] // 先頭の文字
		depth int          // 括弧のネストの深さ
	)

	for _, c := range cs {
		if pref == c {
			return c
		}
	}

	if pref == "(" {
		for {
			if len(clcode) <= 0 {
				return ret
			}
			pref = clcode[:1]
			if pref == "" {
				return ret
			}
			if depth < 0 {
				return ""
			}

			if pref == "(" {
				depth++
				goto endfor
			}
			if pref == ")" {
				depth--
				goto endfor
			}
			if depth == 0 {
				return ret
			}
		endfor:
			ret += pref
			clcode = clcode[1:]
		}
	}

	return pref
}
