
OS_ARCH =\
  darwin-amd64  \
  linux-386     \
  linux-amd64   \
  linux-arm     \
  linux-arm64   \
  windows-386   \
  windows-amd64 \

BUILD_FOLDERS = $(patsubst %,build/%,$(OS_ARCH))

build: $(BUILD_FOLDERS)

$(BUILD_FOLDERS): OS_ARCH=$(patsubst build/%,%,$@)
$(BUILD_FOLDERS): export GOOS=$(word 1,$(subst -, ,$(OS_ARCH)))
$(BUILD_FOLDERS): export GOARCH=$(word 2,$(subst -, ,$(OS_ARCH)))
$(BUILD_FOLDERS): export RELEASE=$(shell git describe)
$(BUILD_FOLDERS):
	@echo "$(GOOS) $(GOARCH)"
	rm -rf $@
	go build -o $@/ ./cmd/...
	cp README.md LICENSE $@/
	cd $@; zip -r ../../dist/gotmpl-$(RELEASE)-$(OS_ARCH).zip *

test:
	go test ./cmd/...

.PHONY: test build $(BUILD_FOLDERS)
