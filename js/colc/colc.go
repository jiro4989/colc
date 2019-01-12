// Package main provides ...
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/jiro4989/colc/combinator/v1"
)

// combinators はコンビネータ定義
var combinators = []combinator.Combinator{
	combinator.Combinator{
		Name:      "S",
		ArgsCount: 3,
		Format:    "{0}{2}({1}{2})",
	},
	combinator.Combinator{
		Name:      "K",
		ArgsCount: 2,
		Format:    "{0}",
	},
	combinator.Combinator{
		Name:      "I",
		ArgsCount: 1,
		Format:    "{0}",
	},
}

func main() {
	calcInput := js.Global.Get("calcInput")
	resultTextArea := js.Global.Get("resultTextArea")

	js.Global.Get("calcButton").Call("addEventListener", "click", func() {
		inputStr := calcInput.Get("value").String()
		combresult := combinator.CalcCLCode(inputStr, combinators, 10)
		resultTextArea.Set("innerHTML", combresult)
	})
}
