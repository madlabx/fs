include build/common.mk
REPO ?= test-images-cluster.xycloud.com/test_registry
TAG ?= test

LDFLAGS += -X "$(MODULE_PATH)/pkg/buildcontext.Version=$(VERSION)" -X "$(MODULE_PATH)/pkg/buildcontext.Commit=$(VERSION_HASH)"
LDFLAGS += -X "$(MODULE_PATH)/pkg/buildcontext.BuildDate=$(DATE)" -X "$(MODULE_PATH)/pkg/buildcontext.Module=$(MODULE)"
LDFLAGS += -X "$(MODULE_PATH)/pkg/buildcontext.Branch=$(BRANCH)"

ENV := test
## Build:
default: prepare  ## Build backend
	$Q CGO_ENABLED=1 $(go) build -ldflags '$(LDFLAGS)' -o $(BIN_INSTALL_DIR)/$(MODULE)  main/$(MODULE).go

prepare:
	@mkdir -p $(IPK_OUTOUT_DIR)
	@mkdir -p $(BIN_INSTALL_DIR)

.PHONY: windows
windows: prepare ## Build windows
	$Q CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(go) build -ldflags '$(LDFLAGS)' -o $(BIN_INSTALL_DIR)/$(MODULE)  main/$(MODULE).go

.PHONY: darwin
darwin: prepare ## Build darwin
	$Q CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(go) build -ldflags '$(LDFLAGS)' -o $(BIN_INSTALL_DIR)/$(MODULE)  main/$(MODULE).go

.PHONY: linux
linux: prepare ## Build linux
	$Q CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(go) build -ldflags '$(LDFLAGS)' -o $(BIN_INSTALL_DIR)/$(MODULE)  main/$(MODULE).go

.PHONY: arm64
arm64: prepare ## Build arm64
	$Q CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 CGO_ENABLED=1 $(go) build -ldflags '$(LDFLAGS)' -o $(BIN_INSTALL_DIR)/$(MODULE)  main/$(MODULE).go

.PHONY: test
test: ## Run base base, exclude cases with skiptest
	$Q $(go) test ./...

.PHONY: test_detail
test-detail: ## with more details
	$Q $(go) test -v ./...

.PHONY: test-all-detail
test-all-detail: ## Run all tests
	$Q $(go) test  --tags="NeedDB" -v  ./...

.PHONY: test-all
test-all: ## Run all tests, include skiptest
	$Q $(go) test  --tags="NeedDB"  ./...

.PHONY: lint
lint: lint-backend lint-commits ## Run all linters

.PHONY: lint-backend
lint-backend: | $(golangci-lint) ## Run backend linters
	$Q $(golangci-lint) run -v

.PHONY: lint-commits
lint-commits: $(commitlint) ## Run commit linters
	$Q ./scripts/commitlint.sh

fmt: $(goimports) ## Format source files
	$Q $(goimports) -local $(MODULE) -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

TAG ?= test

#.PHONY: clean
clean: ## Clean
	@rm -rf $(OUTPUT_DIR)

## Help:
help: ## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target> [options]${RESET}'
	@echo ''
	@echo 'Options:'
	@$(call global_option, "V [0|1]", "enable verbose mode (default:0)")
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

.PHONY: proto

include build/tools.mk
