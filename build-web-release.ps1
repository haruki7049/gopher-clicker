# 1. Build main.wasm
$Env:GOOS = 'js'
$Env:GOARCH = 'wasm'
go build -o main.wasm .
Remove-Item Env:GOOS
Remove-Item Env:GOARCH

# 2. Copy wasm_exec.js from GOROOT to project root
$goroot = go env GOROOT
Copy-Item -Force $goroot/lib/wasm/wasm_exec.js .
