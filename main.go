package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	combinator "github.com/jiro4989/colc/combinator/v1"
	colcio "github.com/jiro4989/colc/io"
)

// options オプション引数
type options struct {
	Version        func() `short:"v" long:"version" description:"バージョン情報"`
	StepCount      int    `short:"s" long:"stepcount" description:"何ステップまで計算するか" default:"-1"`
	OutFile        string `short:"o" long:"outfile" description:"出力ファイルパス"`
	OutFileType    string `short:"t" long:"outfiletype" description:"出力ファイルの種類(なし|json)"`
	Indent         string `short:"i" long:"indent" description:"outfiletypeが有効時に整形して出力する"`
	CombinatorFile string `short:"c" long:"combinatorFile" description:"コンビネータ定義ファイルパス"`
	PrintFlag      bool   `short:"p" long:"print" description:"計算過程を出力する"`
	NoPrintHeader  bool   `short:"n" long:"noprintheader" description:"printフラグON時のヘッダ出力を消す"`
}

type OutValue struct {
	Input   string   `json:"input"`
	Process []string `json:"process"`
	Result  string   `json:"result"`
}
type OutValues []OutValue

// コンビネータ設定
type Combinators []combinator.Combinator

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
	opts, args := parseOptions()

	// コンビネータのファイルパス指定があれば上書き
	if opts.CombinatorFile != "" {
		var err error
		combinators, err = ReadCombinator(opts.CombinatorFile)
		if err != nil {
			panic(err)
		}
	}

	failure := func(err error) {
		panic(err)
	}

	// 引数指定なしの場合は標準入力を処理
	if len(args) < 1 {
		r := os.Stdin
		if err := calcOut(r, opts, out, failure); err != nil {
			panic(err)
		}
		return
	}

	// 引数指定ありの場合はファイル処理
	for _, fn := range args {
		err := colcio.WithOpen(fn, func(r io.Reader) error {
			return calcOut(r, opts, out, failure)
		})
		if err != nil {
			panic(err)
		}
	}
}

// calcOut はCLCodeを計算して、出力する。
// 計算結果を引数の関数に渡し、失敗時は引数に渡した関数を適用する。
func calcOut(r io.Reader, opts options, success func([]string, options) error, failure func(error)) error {
	ss, err := calcCLCode(r, opts)
	if err != nil {
		failure(err)
	}
	return success(ss, opts)
}

// calcCLCode はCLCodeを計算し、スライスで返す。
// OutFileTypeにJSON指定があった場合は、JSON文字列として返す
func calcCLCode(r io.Reader, opts options) ([]string, error) {
	var (
		res []string
		ovs OutValues
		sc  = bufio.NewScanner(r)
	)
	for sc.Scan() {
		line := sc.Text()
		line = strings.Trim(line, " ")

		var (
			s       string
			process []string
		)
		// 出力フラグがある場合は、1ステップ毎に出力
		if opts.PrintFlag {
			// 出力無効化フラグがONなら非表示
			if !opts.NoPrintHeader {
				res = append(res, "=== "+line+" ===")
			}
			var (
				bef = line
				c   = opts.StepCount
			)
			for {
				if c == 0 {
					break
				}
				aft := combinator.CalcCLCode1Time(bef, combinators)
				if bef == aft {
					break
				}
				bef = aft

				switch opts.OutFileType {
				case "json":
					process = append(process, bef)
				}

				res = append(res, bef)
				c--
			}
			s = bef
		} else {
			s = combinator.CalcCLCode(line, combinators, opts.StepCount)
		}

		// JSON出力するときは最後にresを上書きするのでappendしないためにcontinue
		switch opts.OutFileType {
		case "json":
			ov := OutValue{Input: line, Process: process, Result: s}
			ovs = append(ovs, ov)
			continue
		}

		res = append(res, s)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	// JSON出力のときはJSON配列としてすでに完成されたstringをひとつだけ返す
	switch opts.OutFileType {
	case "json":
		var (
			b   []byte
			err error
		)
		if opts.Indent != "" {
			b, err = json.MarshalIndent(ovs, "", opts.Indent)
		} else {
			b, err = json.Marshal(ovs)
		}
		if err != nil {
			return nil, err
		}
		s := string(b)
		res = []string{s}
	}

	return res, nil
}

// out は行配列をオプションに応じて出力する。
// 出力先ファイルが指定されていなければ標準出力する。
// 指定があればファイル出力する。
func out(lines []string, opts options) error {
	if opts.OutFile == "" {
		for _, v := range lines {
			fmt.Println(v)
		}
		return nil
	}

	return colcio.WriteFile(opts.OutFile, lines)
}

// ReadCombinator は指定パスのJSON設定ファイルを読み取る
func ReadCombinator(path string) (Combinators, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var combs Combinators
	if err := json.Unmarshal(b, &combs); err != nil {
		return nil, err
	}
	return combs, nil
}

// parseOptions はコマンドラインオプションを解析する。
// 解析あとはオプションと、残った引数を返す。
func parseOptions() (options, []string) {
	var opts options
	opts.Version = func() {
		fmt.Println(Version)
		os.Exit(0)
	}

	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}

	return opts, args
}
