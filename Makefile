.PHONY: build_win build_mac build_lin win_env mac_env lin_env build

build = CGO_ENABLED=0 \
    GOOS=$(GOOS) \
    GOARCH=$(GOARCH) \
    go build -a \
    -installsuffix cgo \
    -o build/$(GOEXE) \
    main.go

build_win: GOOS=windows
build_win: GOARCH=amd64
build_win: GOEXE=ankr-chain-cli_$(GOOS)_$(GOARCH).exe
build_win:
	@echo "Building win executable"
	@$(build)

build_mac: GOOS=darwin
build_mac: GOARCH=amd64
build_mac: GOEXE=ankr-chain-cli_$(GOOS)_$(GOARCH)
build_mac:
	@echo "Building mac executable"
	@$(build)

build_lin: GOOS=linux
build_lin: GOARCH=amd64
build_lin: GOEXE=ankr-chain-cli_$(GOOS)_$(GOARCH)
build_lin:
	@echo "Building linux executable"
	@$(build)

clean:
	@echo "Cleaning up all the builds"
	@rm -f build/*

build_all: clean build_win build_mac build_lin