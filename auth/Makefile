LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate-auth-api:
	mkdir -p pkg/authV1
	protoc --proto_path api/authV1 \
	--go_out=pkg/authV1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/authV1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
	--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
	api/authV1/auth.proto

generate-access-api:
	mkdir -p pkg/access_v1
	protoc --proto_path api/accessV1 \
          		--go_out=pkg/access_v1 --go_opt=paths=source_relative \
               --plugin=protoc-gen-go=bin/protoc-gen-go \
               --go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative \
               --plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
               api/accessV1/access.proto


vendor-protogen: