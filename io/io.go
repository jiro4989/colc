package io

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// WithOpen はファイルを開き、関数を適用する。
// 自動でファイルをクローズする。
func WithOpen(fn string, f func(r io.Reader) error) error {
	if f == nil {
		return errors.New("適用する関数がnilでした。")
	}
	r, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer r.Close()
	return f(r)
}

// WriteFile はファイル出力する。
// 自動でファイルをクローズする。
func WriteFile(fn string, lines []string) error {
	w, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer w.Close()

	for _, v := range lines {
		fmt.Fprintln(w, v)
	}
	return nil
}
