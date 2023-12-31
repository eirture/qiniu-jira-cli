
#VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)
COMMIT ?= $(shell git rev-parse HEAD)
DATE_FMT = +%Y-%m-%d
BUILD_DATE ?= $(shell date "$(DATE_FMT)")

GO_LDFLAGS := -X github.com/eirture/qiniu-jira-cli/pkg/build.Commit=$(COMMIT) $(GO_LDFLAGS)
GO_LDFLAGS := -X github.com/eirture/qiniu-jira-cli/pkg/build.Version=$(BUILD_DATE) $(GO_LDFLAGS)

.PHONE: build
build:
	@echo "Building..."
	@go build -trimpath -ldflags "${GO_LDFLAGS}" -o bin/ ./cmd/...

.PHONE: install
install:
	@echo "Installing..."
	@go install -trimpath -ldflags "${GO_LDFLAGS}" ./cmd/...

.PHONE: uninstall
uninstall:
	@echo "Uninstalling..."
	@go clean -i ./cmd/...

.PHONE: clean
clean:
	@echo "Cleaning..."
	@rm -rf bin/