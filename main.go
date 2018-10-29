package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	combinator "github.com/jiro4989/colc/combinator/v1"
)

var cs = []combinator.Combinator{
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

// options オプション引数
type options struct {
	Version     func() `short:"v" long:"version" description:"バージョン情報"`
	StepCount   int    `short:"s" long:"stepcount" description:"何ステップまで計算するか"`
	OutFile     string `short:"o" long:"outfile" description:"出力ファイルパス"`
	OutFileType string `short:"t" long:"outfiletype" description:"出力ファイルの種類(なし|json)"`
}

// コンビネータ設定
type Config []CombinatorFormat

type CombinatorFormat struct {
	ArgsCount      int    `json:"argsCount"`
	CombinatorName string `json:"combinatorName"`
	Format         string `json:"format"`
}

// エラー出力ログ
var logger = log.New(os.Stderr, "", 0)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	// オプション引数の解析
	_, args := parseOptions()
	if len(args) < 1 {
		ss, err := calc(os.Stdin)
		if err != nil {
			panic(err)
		}
		for _, s := range ss {
			fmt.Println(s)
		}
	}
}

func calc(r io.Reader) ([]string, error) {
	var res []string
	// 入力をfloatに変換して都度計算
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		line = strings.Trim(line, " ")
		s := combinator.CalcCLCode(line, cs)
		res = append(res, s)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

// ReadConfig は指定パスのJSON設定ファイルを読み取る
func ReadConfig(path string) (Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return config, nil
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
