package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	info := func(msg string) {
		fmt.Println("------------------------------------------------------------")
		fmt.Println("    " + msg)
		fmt.Println("------------------------------------------------------------")
	}

	info("普通の処理をファイル出力する")
	os.Args = []string{
		"main.go",
		"-o",
		"testdata/out/normal_clcode.list",
		"testdata/in/normal_clcode.list",
	}
	main()

	info("普通の処理を標準出力する")
	os.Args = []string{
		"main.go",
		"testdata/in/normal_clcode.list",
	}
	main()

	info("コンビネータ定義ファイルを読み込む")
	os.Args = []string{
		"main.go",
		"-c",
		"config/combinator.json",
		"-o",
		"testdata/out/read_combinator.list",
		"testdata/in/normal_clcode.list",
	}
	main()

	info("計算過程を標準出力する")
	os.Args = []string{
		"main.go",
		"-p",
		"testdata/in/normal_clcode.list",
	}
	main()

	info("計算過程を標準出力するが、ヘッダを出力しない")
	os.Args = []string{
		"main.go",
		"-pn",
		"testdata/in/normal_clcode.list",
	}
	main()

	info("計算結果をJSON出力する")
	os.Args = []string{
		"main.go",
		"-t",
		"json",
		"testdata/in/normal_clcode.list",
	}
	main()

	info("計算結果をJSON出力する + インデント")
	os.Args = []string{
		"main.go",
		"-t",
		"json",
		"-i",
		"  ",
		"testdata/in/normal_clcode.list",
	}
	main()

	info("計算結果をJSON出力する + TABインデント")
	os.Args = []string{
		"main.go",
		"-t",
		"json",
		"-i",
		"\t",
		"testdata/in/normal_clcode.list",
	}
	main()

	info("計算結果をJSON出力する + インデント + 計算過程出力")
	os.Args = []string{
		"main.go",
		"-pt",
		"json",
		"-i",
		"  ",
		"testdata/in/normal_clcode.list",
	}
	main()

}

func TestCalcOut(t *testing.T) {
	f := func(ss ...string) io.Reader {
		return bytes.NewBufferString(strings.Join(ss, "\n"))
	}
	o1 := options{StepCount: -1}
	o2 := options{StepCount: 1}
	o3 := options{StepCount: -1, CombinatorFile: "config/combinator.json"}
	type TD struct {
		r       io.Reader
		opts    options
		success func([]string, options) error
		failure func(error)
	}
	tds := []TD{
		TD{
			r:    f("Sxyz"),
			opts: o1,
			success: func(ss []string, opts options) error {
				assert.Equal(t, []string{"xz(yz)"}, ss)
				return nil
			},
			failure: func(err error) {
				assert.NoError(t, err)
			},
		},
		TD{
			r:    f("Sxyz", "(SSSS)"),
			opts: o1,
			success: func(ss []string, opts options) error {
				assert.Equal(t, []string{"xz(yz)", "SS(SS)"}, ss)
				return nil
			},
			failure: func(err error) {
				assert.NoError(t, err)
			},
		},
		TD{
			r:    f("KKxy"),
			opts: o2,
			success: func(ss []string, opts options) error {
				assert.Equal(t, []string{"Ky"}, ss)
				return nil
			},
			failure: func(err error) {
				assert.NoError(t, err)
			},
		},
		TD{
			r:    f("<true>xy"),
			opts: o3,
			success: func(ss []string, opts options) error {
				assert.Equal(t, []string{"x"}, ss)
				return nil
			},
			failure: func(err error) {
				assert.NoError(t, err)
			},
		},
		TD{
			r:    f("SBKI"),
			opts: o3,
			success: func(ss []string, opts options) error {
				assert.Equal(t, []string{"BI(KI)"}, ss)
				return nil
			},
			failure: func(err error) {
				assert.NoError(t, err)
			},
		},
	}
	for _, v := range tds {
		r, opts, success, failure := v.r, v.opts, v.success, v.failure
		err := calcOut(r, opts, success, failure)
		assert.NoError(t, err)
	}
}

func TestCalcCLCode(t *testing.T) {
	f := func(ss ...string) io.Reader {
		return bytes.NewBufferString(strings.Join(ss, "\n"))
	}
	type TD struct {
		r    io.Reader
		opts options
		s    []string
		desc string
	}

	o1 := options{StepCount: -1}
	o2 := options{StepCount: 1, Indent: "  "}
	o3 := options{StepCount: 0}
	o4 := options{StepCount: -1, OutFileType: "json"}
	o5 := options{StepCount: -1, OutFileType: "json", PrintFlag: true}
	o6 := options{StepCount: -1, OutFileType: "json", PrintFlag: true, Indent: "  "}
	o7 := options{StepCount: -1, PrintFlag: true}

	tds := []TD{
		TD{
			r:    f("Sxyz"),
			opts: o1,
			s:    []string{"xz(yz)"},
			desc: "正常系:オプションなし",
		},
		TD{
			r:    f("Sxyz", "SKII"),
			opts: o1,
			s:    []string{"xz(yz)", "I"},
			desc: "正常系:複数行処理",
		},
		TD{
			r:    f("Sxyz", "SKII"),
			opts: o2,
			s:    []string{"xz(yz)", "KI(II)"},
			desc: "正常系:計算回数指定。及び、JSON指定がない状態でインデント指定をしても意味がない",
		},
		TD{
			r:    f("Sxyz", "SKII"),
			opts: o3,
			s:    []string{"Sxyz", "SKII"},
			desc: "正常系:計算しない",
		},
		TD{
			r:    f("SSSSS"),
			opts: o1,
			s:    []string{"SS((SS)S)"},
			desc: "正常系:ネスト括弧の計算をする",
		},
		TD{
			r:    f("Sxyz"),
			opts: o4,
			s:    []string{`[{"input":"Sxyz","process":null,"result":"xz(yz)"}]`},
			desc: "正常系:計算結果のJSON出力",
		},
		TD{
			r:    f("Sxyz", "Sxyz"),
			opts: o4,
			s:    []string{`[{"input":"Sxyz","process":null,"result":"xz(yz)"},{"input":"Sxyz","process":null,"result":"xz(yz)"}]`},
			desc: "正常系:計算結果の複数JSON出力",
		},
		TD{
			r:    f("SKIx"),
			opts: o5,
			s:    []string{`[{"input":"SKIx","process":["Kx(Ix)","x"],"result":"x"}]`},
			desc: "正常系:計算結果のJSON出力",
		},
		TD{
			r:    f("Sxyz", "Sxyz"),
			opts: o6,
			s: []string{`[
  {
    "input": "Sxyz",
    "process": [
      "xz(yz)"
    ],
    "result": "xz(yz)"
  },
  {
    "input": "Sxyz",
    "process": [
      "xz(yz)"
    ],
    "result": "xz(yz)"
  }
]`},
			desc: "正常系:2スペースインデントされたJSON出力",
		},
		TD{
			r:    f("SKIx"),
			opts: o7,
			s:    []string{"=== SKIx ===", "Kx(Ix)", "x", "x"},
			desc: "正常系:計算過程の出力",
		},
	}
	for _, v := range tds {
		r, opts, expect, desc := v.r, v.opts, v.s, v.desc
		actual, err := calcCLCode(r, opts)
		assert.Equal(t, expect, actual, desc, r, opts)
		assert.NoError(t, err, desc)
	}
}
