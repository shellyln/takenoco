

DEVNUL     := /dev/null
DIRSEP     := /
SEP        := :
RM_F       := rm -f
RM_RF      := rm -rf
CP         := cp
CP_FLAGS   :=
CP_R       := cp -RT
CP_R_FLAGS :=
FINDFILE   := find . -type f -name
WHICH      := which
PRINTENV   := printenv
GOOS       := linux
GOARCH     := amd64
GOARM      :=
GOXOS      := darwin windows linux
GOXARCH    := 386 amd64 arm
GOXARM     := 7
GOCMD      := go
GOBUILD    := $(GOCMD) build
GOTIDY     := $(GOCMD) mod tidy
GOCLEAN    := $(GOCMD) clean
GOTEST     := $(GOCMD) test
GOVET      := $(GOCMD) vet
GOLINT     := $(GOPATH)/bin/golint -set_exit_status
TINYGOCMD  := tinygo
SRCS       :=
TARGET_CLI := ./
BIN_CLI    := app


ifeq ($(OS),Windows_NT)
    BIN_CLI := $(BIN_CLI).exe
    ifeq ($(MSYSTEM),)
        SHELL      := cmd.exe
        DEVNUL     := NUL
        DIRSEP     := \\
        SEP        := ;
        RM_F       := del /Q
        RM_RF      := rmdir /S /Q
        CP         := copy
        CP_FLAGS   := /Y
        CP_R       := xcopy
        CP_R_FLAGS := /E /I /Y
        FINDFILE   := cmd.exe /C 'where /r . '
        WHICH      := where
        PRINTENV   := set
    endif
endif

define normalize_dirsep
    $(subst /,$(DIRSEP),$1)
endef

define find_file
    $(subst $(subst \,/,$(CURDIR)),.,$(subst \,/,$(shell $(FINDFILE) $1)))
endef


# Usage of cp -R and cp
# $(CP_R) $(call normalize_dirsep,path/to/src) $(call normalize_dirsep,path/to/dest) $(CP_R_FLAGS)
# $(CP)   $(call normalize_dirsep,path/to/src) $(call normalize_dirsep,path/to/dest) $(CP_FLAGS)


SRCS     := $(call find_file,"*.go")
VERSION  := $(shell git describe --tags --abbrev=0 2> $(DEVNUL) || echo "0.0.0-alpha.1")
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS  := -ldflags="-s -w -buildid= -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""


.PHONY: printenv clean tidy test lint
all: clean test build


printenv:
	@echo SHELL      : $(SHELL)
	@echo CURDIR     : $(CURDIR)
	@echo DEVNUL     : $(DEVNUL)
	@echo DIRSEP     : "$(DIRSEP)"
	@echo SEP        : "$(SEP)"
	@echo WHICH GO   : $(shell $(WHICH) $(GOCMD))
	@echo GOOS       : $(GOOS)
	@echo GOARCH     : $(GOARCH)
	@echo GOARM      : $(GOARM)
	@echo VERSION    : $(VERSION)
	@echo REVISION   : $(REVISION)
	@echo SRCS       : $(SRCS)
	@echo LDFLAGS    : $(LDFLAGS)
	@echo TARGET_CLI : $(TARGET_CLI)
	@echo BIN_CLI    : $(BIN_CLI)


clean:
	$(GOCLEAN)
	-$(RM_F) $(BIN_CLI)

cleantest:
	$(GOCLEAN) -testcache


tidy:
	$(GOTIDY)


test:
	$(GOTEST) ./...

testinfo:
	$(GOTEST) -gcflags=-m ./...

cover:
	$(GOTEST) -cover ./...

lint:
	@echo "Run go vet..."
	$(GOVET) ./...


$(BIN_CLI): export CGO_ENABLED:=0
$(BIN_CLI): $(SRCS)
	$(GOBUILD) \
	    -a -tags osusergo,netgo -installsuffix netgo \
	    -trimpath \
	    $(LDFLAGS) \
	    -o $(BIN_CLI) $(TARGET_CLI)


$(BIN_CLI)_quick: $(SRCS)
	$(GOBUILD) -o $(BIN_CLI) $(TARGET_CLI)


$(BIN_CLI)_info: $(SRCS)
	$(GOBUILD) -gcflags=-m -o $(BIN_CLI) $(TARGET_CLI)


build: $(BIN_CLI) ;


quickbuild: $(BIN_CLI)_quick ;

buildinfo: $(BIN_CLI)_info ;


xbuild: export GOOS:=$(GOOS)
xbuild: export GOARCH:=$(GOARCH)
xbuild: export GOARM:=$(GOARM)
xbuild: build ;


wasm: export GOOS:=js
wasm: export GOARCH:=wasm
wasm: export GOCMD:=go
wasm:
	$(TINYGOCMD) build -tags wasm -o web/go.wasm ./cmd/wasm


fatwasm: export GOOS:=js
fatwasm: export GOARCH:=wasm
fatwasm: export GOCMD:=go
fatwasm:
	$(GOCMD) build -tags wasm -o web/go.wasm ./cmd/wasm


docker:
	docker build -t shellyln/takenoco:$(VERSION) .


docker-test:
	docker build -t shellyln/takenoco:rev-$(REVISION) .


doc:
	godoc -http=:6060
