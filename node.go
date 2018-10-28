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

// CalcCLCode は計算不可能になるまで計算した結果を返す。
func CalcCLCode(clcode string, cs []Combinator) string {
	// コンビネータリストから名前だけのコンビネータを生成
	var cns []string
	for _, c := range cs {
		cns = append(cns, c.Name)
	}

	// break判定用。計算前と後で一致していたらbreak
	var pref string
	bef := clcode
	for {
		pref = getPrefixCombinator(clcode, cns)
		clcode = clcode[len(pref):]
		// コンビネータが括弧で括られていたら、括弧を除去
		if strings.HasPrefix(pref, "(") {
			pref = pref[1 : len(pref)-1]
		}
		if pref == "" {
			break
		}

		// 先頭コンビネータが定義済みコンビネータの中にあればセット
		var co Combinator
		for _, c := range cs {
			if c.Name == pref {
				co = c
				break
			}
		}

		// 計算前のデータを保存。break判定用
		bef = clcode

		// 定義済みコンビネータが必要とする分コンビネータを取得
		var args []string
		for i := 0; i < co.ArgsCount; i++ {
			c := getPrefixCombinator(clcode, cns)
			args = append(args, c)
			clcode = clcode[len(c):]
		}
		// 計算結果は先頭コンビネータ分だけなので、計算されなかった分のコンビネ
		// ータと結合
		clcode = CalcHeadCombinator(args, co) + clcode

		// 計算前と後が同じ == 計算不可能な状態になったら終了
		if bef == clcode {
			break
		}
	}
	return pref + clcode
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
	var (
		co    Combinator
		found bool
	)
	for _, c := range cs {
		if c.Name == pref {
			co = c
			found = true
			break
		}
	}

	// マッチするコンビネータがない場合は空配列を返す
	if !found {
		return []string{}
	}

	clcode = clcode[len(co.Name):]
	var args []string
	for i := 0; i < co.ArgsCount; i++ {
		c := getPrefixCombinator(clcode, cns)
		// 引数よりも見つかったコンビネータ数が少ないときは終了
		if c == "" {
			return []string{}
		}
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

	// 先頭のが定義済みコンビネータだったら返却
	for _, c := range cs {
		if strings.HasPrefix(clcode, c) {
			return c
		}
	}

	var (
		ret   string
		pref  = clcode[:1] // 先頭の文字
		depth int          // 括弧のネストの深さ
	)

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
		endfor:
			ret += pref
			clcode = clcode[1:]

			if depth == 0 {
				return ret
			}
		}
	}

	return pref
}
