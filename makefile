BIN_FOLDER=bin
BIN_LINUX=gotunnel_linux
BIN_DARWIN=gotunnel_darwin
BIN_WINDOWS=gotunnel_windows.exe

all: clean linux_build darwin_build windows_build

linux: clean linux_build
darwin: clean darwin_build
windows: clean windows_build

clean:
	rm -rf $(BIN_FOLDER)


linux_build:
	CGO_ENABLED=0 GOOS=linux go build -o $(BIN_FOLDER)/$(BIN_LINUX)


darwin_build:
	CGO_ENABLED=0 GOOS=darwin go build -o $(BIN_FOLDER)/$(BIN_DARWIN)


windows_build:
	CGO_ENABLED=0 GOOS=windows go build -o $(BIN_FOLDER)/$(BIN_WINDOWS)

