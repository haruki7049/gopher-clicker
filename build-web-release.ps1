# 1. Build main.wasm
$Env:GOOS = 'js'
$Env:GOARCH = 'wasm'
## 1.1. Create a public dir
New-Item -ItemType Directory public
## 1.2. Build!!
go build -o public/main.wasm .
## 1.3. Clean up environment variables
Remove-Item Env:GOOS
Remove-Item Env:GOARCH

# 2. Copy wasm_exec.js from GOROOT to public dir
$goroot = go env GOROOT
Copy-Item -Force $goroot/lib/wasm/wasm_exec.js ./public

# 3. Copy index.html from project root to public dir
Copy-Item -Force ./index.html ./public
