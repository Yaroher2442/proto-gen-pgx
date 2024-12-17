RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(PATH):$(LOCAL_BIN)

default: help

.PHONY: help
help: # Показывает информацию о каждом рецепте в Makefile
	@grep -E '^[a-zA-Z0-9 _-]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: .bin_deps
.bin_deps: # Устанавливает зависимости необходимые для работы приложения
	$(info Installing binary dependencies...)
	mkdir -p $(LOCAL_BIN)

.PHONY: .install_linter
.install_linter: # Устанавливает линтер
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Downloading golangci-lint latest)
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
endif

.PHONY: install
install: .install_linter .bin_deps # Устанавливает все зависимости для работы приложения

.PHONY: tests
tests: # Запускает юнит тесты с ковереджем
	go test -race -coverprofile=coverage.out ./...

.PHONY: linter
linter: # Запуск линтеров
	$(LOCAL_BIN)/golangci-lint cache clean && \
	$(LOCAL_BIN)/golangci-lint run

.PHONY: linter_fix
linter_fix: # Запуск линтеров с фиксом где возможно
	$(LOCAL_BIN)/golangci-lint cache clean && \
	$(LOCAL_BIN)/golangci-lint run --fix

.PHONY: .generate-descriptor
.generate-descriptor:
	protoc --go_out=./internal/domain/descriptorpb --go_opt=paths=source_relative --proto_path=./proto ./proto/descriptor.proto

.PHONY: build
build:
	go build -o cmd/protoc-gen-gopgx ./cmd/protoc-gen-gopgx

examples_files = $(shell find "./examples" -type f -name '*.proto')
.PHONY: .run_example
.run_example: # Запуск примера
	$(MAKE) build
	rm -rf ./examples/basic/basicpb ./examples/basic/proto/descriptor.proto
	cp -r ./proto/descriptor.proto ./examples/basic/proto/descriptor.proto
	mkdir -p ./examples/basic/basicpb
	protoc --go_out=./examples/basic/basicpb \
	--go_opt=paths=source_relative \
	--plugin=./cmd/protoc-gen-gopgx/protoc-gen-gopgx \
	--gopgx_out=./examples/basic/basicpb \
	--proto_path=./examples/basic/proto $(examples_files)

branch=main
.PHONY: revision
revision: # Создание тега
	@if [ -e $(tag) ]; then \
		echo "error: Specify version 'tag='"; \
		exit 1; \
	fi
	git tag -d ${tag} || true
	git push --delete origin ${tag} || true
	git tag $(tag)
	git push origin $(tag)
