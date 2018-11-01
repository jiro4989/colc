package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	o2 := options{StepCount: 1}
	o3 := options{StepCount: 0}

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
			desc: "正常系:計算回数指定",
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
	}
	for _, v := range tds {
		r, opts, expect, desc := v.r, v.opts, v.s, v.desc
		actual, err := calcCLCode(r, opts)
		assert.Equal(t, expect, actual, desc, r, opts)
		assert.NoError(t, err, desc)
	}
}
