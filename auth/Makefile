LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-auth-api
	make generate-access-api

generate-auth-api:
	mkdir -p pkg/authV1
	protoc --proto_path api/authV1 --proto_path vendor.protogen \
	--go_out=pkg/authV1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/authV1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/authV1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	api/authV1/auth.proto

generate-access-api:
	mkdir -p pkg/access_v1
	protoc --proto_path api/accessV1 --proto_path vendor.protogen \
           --go_out=pkg/access_v1 --go_opt=paths=source_relative \
           --plugin=protoc-gen-go=bin/protoc-gen-go \
           --go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative \
           --plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
           --validate_out lang=go:pkg/access_v1 --validate_opt=paths=source_relative \
           --plugin=protoc-gen-validate=bin/protoc-gen-validate \
           api/accessV1/access.proto


vendor-proto:
	@if [ ! -d vendor.protogen/google ]; then \
      		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
      		mkdir -p vendor.protogen/google/ &&\
      		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
      		rm -rf google.protogen/googleapis ;\
    fi
	@if [ ! -d vendor.protogen/validate ]; then \
    		mkdir vendor.protogen/validate && \
    		git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate && \
    		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate && \
    		rm -rf vendor.protogen/protoc-gen-validate ;\
    fi
	@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
        	mkdir -p vendor.protogen/protoc-gen-openapiv2/options && \
            git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 && \
            mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
            rm -rf vendor.protogen/openapiv2 ;\
    fi