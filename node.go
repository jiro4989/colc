package main

import (
	"errors"
	"fmt"
	"strings"
)

// Node はコンビネータの名前と、その子ノードを持つ。
type Node struct {
	Name  string
	Nodes []Node
}

// ParseCLCode はclcodeを解析してNode配列を返す。
func ParseCLCode(clcode string, config ...Config) ([]Node, error) {
	// 空白文字は無視するので削除
	clcode = strings.Replace(clcode, " ", "", -1)
	clcode = strings.Replace(clcode, "　", "", -1)
	clcode = strings.Replace(clcode, "\t", "", -1)

	if clcode == "" {
		return nil, errors.New("clcodeが空文字不正")
	}

	// コンフィグ指定があればコンビネータを設定する
	// コンフィグ指定がなければソース内のコンビネータを使用する。
	var combinators []string
	if config == nil {
		combinators = []string{"S", "K", "I"}
	} else {
		for _, c := range config[0] {
			combinators = append(combinators, c.CombinatorName)
		}
	}

	var (
		nodes []Node                                     // 返却する値
		node  Node                                       // ループ内で使う一時node
		c     = getPrefixCombinator(clcode, combinators) // 先頭の文字
		depth = 0                                        // 括弧のネストの深さ
	)
	clcode = clcode[len(c):]
	for {
		if depth < 0 {
			msg := fmt.Sprintf("()の対応関係が不正 depth=%d", depth)
			err := errors.New(msg)
			return nil, err
		}
		if c == "" {
			break
		}

		for _, combinator := range combinators {
			if c == combinator {
				goto endfor
			}
		}
		if c == "(" {
			depth++
			goto endfor
		}
		if c == ")" {
			depth--
			goto endfor
		}
	endfor:
		node.Name += c
		// 括弧のネストが0の時にNodeをスライスに追加する
		if depth == 0 {
			if node.Name == "()" {
				return nil, errors.New("()のみのデータは不正")
			}

			// 括弧でくくられてるやつは展開する
			nm := node.Name
			if strings.HasPrefix(nm, "(") {
				nm = nm[1:]
				nm = nm[:len(nm)-1]
				var err error
				node.Nodes, err = ParseCLCode(nm)
				if err != nil {
					return nodes, err
				}
			}

			nodes = append(nodes, node)
			node = Node{}
		}
		if 0 < len(clcode) {
			c = getPrefixCombinator(clcode, combinators)
			clcode = clcode[len(c):]
			continue
		}
		c = ""
	}
	if depth != 0 {
		msg := fmt.Sprintf("()の対応関係が不正 depth=%d", depth)
		err := errors.New(msg)
		return nil, err
	}
	return nodes, nil
}

// getPrefixCombinator はCLCodeの先頭のコンビネータを返す。
// 引数に渡している定義済みコンビネータが存在した場合、複数文字でも返す。
// 定義済みコンビネータにマッチしない場合は、先頭1文字を返す。
func getPrefixCombinator(clcode string, cs []string) string {
	if len(clcode) < 1 {
		return ""
	}
	for _, c := range cs {
		if strings.HasPrefix(clcode, c) {
			return c
		}
	}
	return clcode[:1]
}
