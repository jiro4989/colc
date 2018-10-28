APPNAME := $(shell basename `pwd`)
VERSION := v$(shell gobump show -r)
SRCS := $(shell find . -name "*.go" -type f )
LDFLAGS := -ldflags="-s -w \
	-extldflags \"-static\""
XBUILD_TARGETS := \
	-os="windows linux darwin" \
	-arch="386 amd64" 
DIST_DIR := dist/$(VERSION)
README := README.md
EXTERNAL_TOOLS := \
	github.com/golang/dep/cmd/dep \
	github.com/mitchellh/gox \
	github.com/tcnksm/ghr \
	github.com/motemen/gobump/cmd/gobump \
	github.com/alecthomas/gometalinter

help: ## ドキュメントのヘルプを表示する。
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: $(SRCS) ## ビルド
	go build $(LDFLAGS) -o bin/$(APPNAME) .

install: build ## インストール
	go install

.PHONY: xbuild
xbuild: $(SRCS) bootstrap ## クロスコンパイル
	gox $(LDFLAGS) $(XBUILD_TARGETS) --output "$(DIST_DIR)/{{.Dir}}_{{.OS}}_{{.Arch}}/{{.Dir}}"

.PHONY: archive
archive: xbuild ## クロスコンパイルしたバイナリとREADMEを圧縮する
	find $(DIST_DIR)/ -mindepth 1 -maxdepth 1 -a -type d \
		| while read -r d; \
		do \
			cp $(README) $$d/ ; \
		done
	cd $(DIST_DIR) && \
		find . -maxdepth 1 -mindepth 1 -a -type d  \
		| while read -r d; \
		do \
			tar czf $$d.tar.gz $$d; \
		done

.PHONY: release
release: bootstrap test archive ## GitHubにリリースする
	ghr $(VERSION) $(DIST_DIR)/

.PHONY: test
test: ## テストコードを実行する
	-gometalinter
	go test -cover ./...

.PHONY: clean
clean: ## バイナリ、配布物ディレクトリを削除する
	-rm -rf bin
	-rm -rf $(DIST_DIR)

.PHONY: deps
deps: bootstrap ## 依存ライブラリを更新する
	dep ensure

.PHONY: bootstrap
bootstrap: ## 外部ツールをインストールする
	for t in $(EXTERNAL_TOOLS); do \
		echo "Installing $$t ..." ; \
		go get $$t ; \
	done
	gometalinter --install --update

graph: ## グラフ画像を生成する
	docker build ./graphviz -t graphviz
	docker run -v `pwd`/doc:/root/doc -v `pwd`/script/generate_graph.sh:/generate_graph.sh -it graphviz /generate_graph.sh
