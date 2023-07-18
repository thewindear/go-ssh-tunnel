# Makefile

# 交叉编译为MacOS
build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/ssh-tunnel-app

# 交叉编译为Linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/ssh-tunnel-app

# 交叉编译为Windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/ssh-tunnel-app.exe

# 默认目标，执行全部编译
all: build-macos build-linux build-windows