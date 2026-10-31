#!/bin/sh
# 把画廊编译成 WebAssembly，放到 public/preview/，文档页的预览和在线画廊都加载它。
set -e
cd "$(dirname "$0")"
out=public/preview
GOOS=js GOARCH=wasm CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o "$out/gallery.wasm" ../cmd/gallery
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$out/wasm_exec.js"
