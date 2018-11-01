package combinator

import (
	"fmt"
	"strings"
)

// Combinator はコンビネータである。
type Combinator struct {
	Name      string
	ArgsCount int
	Format    string
}

// CalcCLCode は計算不可能になるまで計算した結果を返す。
func CalcCLCode(clcode string, cs []Combinator, n int) string {
	ret := CalcCLCode1Time(clcode, cs)
	if ret == clcode {
		return ret
	}
	return CalcCLCode(ret, cs, n)
}

// CalcCLCode1Time は先頭のコンビネータを一度だけ計算する。
// 括弧があっても展開して1回計算する。
func CalcCLCode1Time(clcode string, cs []Combinator) string {
	pref, args, suff := splitPrefixArgsSuffixCombinators(clcode, cs)

	// 括弧が出現したときは括弧を展開して計算
	if strings.HasPrefix(pref, "(") && strings.HasSuffix(pref, ")") {
		suff := clcode[len(pref):]
		pref = trimBracket(pref)
		return CalcCLCode1Time(pref+suff, cs)
	}

	// 先頭コンビネータが定義済みコンビネータの中にあればセット
	var co Combinator
	var found = false
	for _, c := range cs {
		if c.Name == pref {
			co = c
			found = true
			break
		}
	}
	if !found || len(args) != co.ArgsCount {
		return clcode
	}
	return calcCombinatorArgs(args, co) + suff
}

// trimBracket は括弧で括られたCLCodeから括弧を除く。
// 除く対象は、複数のコンビネータを一つのコンビネータとしてラッピングしてしまっ
// ている一番外に1つ以上つづく括弧だけである。
func trimBracket(s string) string {
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		return trimBracket(s[1 : len(s)-1])
	}
	return s
}

// calcCombinatorArgs は引数コンビネータをコンビネータで計算する。
func calcCombinatorArgs(cs []string, co Combinator) string {
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

// splitPrefixArgsSuffixCombinators はCLCodeを先頭コンビネータ、引数コンビネータ
// 、計算対象外コンビネータに分けて返却する。
func splitPrefixArgsSuffixCombinators(clcode string, cs []Combinator) (string, []string, string) {
	pref := getPrefixCombinator(clcode, cs)
	args := getCombinatorArgs(clcode, cs)
	suff := clcode[len(pref+strings.Join(args, "")):]
	return pref, args, suff
}

// getPrefixCombinator はCLCodeの先頭のコンビネータを返す。
// 引数に渡している定義済みコンビネータが存在した場合、複数文字でも返す。
func getPrefixCombinator(clcode string, cs []Combinator) string {
	if len(clcode) < 1 {
		return ""
	}

	// 先頭のが定義済みコンビネータだったら返却
	for _, c := range cs {
		nm := c.Name
		if strings.HasPrefix(clcode, nm) {
			return nm
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

// getCombinatorArgs は先頭のコンビネータを判定し、計算対象のコンビネータを返却する。
func getCombinatorArgs(clcode string, cs []Combinator) []string {
	pref := getPrefixCombinator(clcode, cs)

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
		c := getPrefixCombinator(clcode, cs)
		// 引数よりも見つかったコンビネータ数が少ないときは終了
		if c == "" {
			return []string{}
		}
		args = append(args, c)
		clcode = clcode[len(c):]
	}

	return args
}
