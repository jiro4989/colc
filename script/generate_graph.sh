#!/bin/bash

set -eu

cd /root/doc

for path in graphviz/*.dot; do
	fn=$(basename $path)
	fn=$(echo "${fn%%.*}")
	dot -Tpng "$path" -o "graphviz/$fn.png"
done
