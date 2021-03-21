PROJECTNAME=$(shell basename "$(PWD)")

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

BUILDOPTS = "-tags=jsoniter -trimpath -ldflags=\"-s -w\""

all: build_linux_amd64 build_linux_arm64 build_mac_amd64 build_windows_amd64

go_mod_tidy: 
	echo "- go mod tidy"
	go mod tidy

swag: 
	echo "- swag"
	swag init --parseDependency --parseInternal

build_linux_amd64: go_mod_tidy swag
	echo "- build_linux_amd64"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${BUILDOPTS} -o bin/${PROJECTNAME}-linux-amd64 main.go

build_linux_arm64: go_mod_tidy swag
	echo "- build_linux_arm64"
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ${BUILDOPTS} -o bin/${PROJECTNAME}-linux-arm64 main.go

build_mac_amd64: go_mod_tidy swag
	echo "- build_mac_amd64"
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${BUILDOPTS} -o bin/${PROJECTNAME}-mac-amd64 main.go

build_windows_amd64: go_mod_tidy swag
	echo "- build_windows_amd64"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${BUILDOPTS} -o bin/${PROJECTNAME}-windows-amd64.exe main.go
