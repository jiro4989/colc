package combinator

import (
	"strings"
)

// Combinator はコンビネータである。
type Combinator struct {
	Name      string `json:"name"`
	ArgsCount int    `json:"argsCount"`
	Format    string `json:"format"`
}
type Combinators []Combinator

type Node struct {
	Text  string
	Nodes []*Node
}

func (n *Node) IsLeaf() bool {
	return len(n.Nodes) <= 0
}

func Parse(code string, cs Combinators) Node {
	node := Node{Text: code}
	for code != "" {
		first := getPrefixCombinator(code, cs)
		first = trimBracket(first)
		code = code[len(first):]

		n := Node{Text: first}
		node.Nodes = append(node.Nodes, &n)
	}
	// if !node.IsLeaf() {
	// 	for i, v := range node.Nodes {
	// 		vv := *v
	// 		n := Parse(vv.Text, cs)
	// 		node.Nodes[i] = &n
	// 	}
	// }
	return node
}

// getPrefixCombinator はCLCodeの先頭のコンビネータを返す。
// 引数に渡している定義済みコンビネータが存在した場合、複数文字でも返す。
func getPrefixCombinator(clcode string, cs Combinators) string {
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

// trimBracket は括弧で括られたCLCodeから括弧を除く。
// 除く対象は、複数のコンビネータを一つのコンビネータとしてラッピングしてしまっ
// ている一番外に1つ以上つづく括弧だけである。
func trimBracket(s string) string {
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		return trimBracket(s[1 : len(s)-1])
	}
	return s
}
