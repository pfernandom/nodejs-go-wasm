npm install --save express esbuild
go mod init wasm-test.com

GOOS=js GOARCH=wasm go build -o json.wasm

go get "syscall/js"

```
 "go.toolsEnvVars": {
    "GOOS": "js",
    "GOARCH": "wasm"
  }
```
