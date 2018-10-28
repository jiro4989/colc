# colc (COmbinator Logic Calculator)
[![Build Status](https://travis-ci.org/jiro4989/colc.svg?branch=master)](https://travis-ci.org/jiro4989/colc)

Combinator Logicをコマンドラインから使うためのツール

## コンビネータ論理(Combinator Logic)とは
引数に関数を受け取る関数(コンビネータ)のみで計算をするという考え方をコンビネータ
論理という。
チューリング完全であることが証明されているため、コンピュータで可能な計算は全てコ
ンビネータだけで計算が可能である。

以下にコンビネータ論理の主要な関数3つを定義。

```
# Sabc -> ac(bc)
S, 3, 02(12)

# Kab
I, 2, 0

# Ia
I, 1, 0
```

### コンビネータ
![Sコンビネータとコンビネータの分割](doc/graphviz/s_combinator.png)

![SKIの計算の流れ](doc/graphviz/mix_combinator.png)

## 使い方

```bash
echo "Sxyz" | colc
# -> xz(yz)

colc clcode.txt

colc -f clcode1.txt -f clcode2.txt

# ファイル出力
colc clcode.txt -o out.txt

# JSON出力
colc clcode.txt -t json
colc clcode.txt -t json -o out.json
```

## 開発
### ヘルプ

```bash
make
```

### バイナリの生成

```bash
make build
```

### グラフ画像の生成

```bash
make graph
```
