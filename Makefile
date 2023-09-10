.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest

.PHONY: config
# generate internal proto
config:
	find app -maxdepth 1 -mindepth 1 -type d -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) config'

.PHONY: api
# generate api proto
api:
	find app -maxdepth 1 -mindepth 1 -type d -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) api'

.PHONY: errors
# generate errors proto
errors:
	find app -maxdepth 1 -mindepth 1 -type d -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) errors'

.PHONY: validate
# generate validate proto
validate:
	find app -maxdepth 1 -mindepth 1 -type d -print | xargs -L 1 bash -c 'cd "$$0" && pwd && $(MAKE) validate'

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: service
# create configs
service:
	kratos new app/$(name) --nomod
	kratos proto add api/$(name)/v1/$(name).proto && \
	kratos proto client api/$(name)/v1/$(name).proto && \
	kratos proto server api/$(name)/v1/$(name).proto -t app/$(name)/internal/service && \
	cd app/$(name) && echo "include ../../app_makefile" > ./Makefile && cd ../../

.PHONY: all
# generate all
all:
	make api;
	make config;
	make generate;