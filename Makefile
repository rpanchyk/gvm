.PHONY: build deps build-linux-amd64 build-macosx-amd64 build-macosx-arm64 build-windows-amd64

define make_build
	rm -f builds/tmp/*
	GOOS=$(1) GOARCH=$(2) go build -o builds/tmp/.gvm/
	cp -f config.toml builds/tmp/.gvm/
	cd builds/tmp && zip --recurse-paths --move ../$(basename $3)-$(4)-$(2).zip . && cd -
endef

# Batch build
build: deps build-linux-amd64 build-macosx-amd64 build-macosx-arm64 build-windows-amd64

# Dependencies
deps:
	go mod tidy && go mod vendor

# Linux
build-linux-amd64:
	$(call make_build,linux,amd64,gvm,linux)

# MacOSX
build-macosx-amd64:
	$(call make_build,darwin,amd64,gvm,macosx)

build-macosx-arm64:
	$(call make_build,darwin,arm64,gvm,macosx)

# Windows
build-windows-amd64:
	$(call make_build,windows,amd64,gvm.exe,windows)
