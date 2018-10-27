package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
)

// options オプション引数
type options struct {
	Version     func()   `short:"v" long:"version" description:"バージョン情報"`
	CLCode      []string `short:"c" long:"clcode" description:"計算対象のCLCode"`
	InFile      string   `short:"f" long:"infile" description:"計算対象の書かれたテキストファイル"`
	OutFile     string   `short:"o" long:"outfile" description:"出力ファイルパス"`
	OutFileType string   `short:"t" long:"outfiletype" description:"出力ファイルの種類(なし|json)"`
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
	opts, args := parseOptions()
	fmt.Println(opts, args) // DEBUG

	config, _ := ReadConfig("config/combinator.json")
	v, _ := ParseCLCode("Sxyz", config)
	fmt.Println(v)

	v, _ = ParseCLCode("S(SKI)", config)
	fmt.Println(v)
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
