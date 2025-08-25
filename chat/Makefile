docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/dimastephen/chatserver:v0.1 .
	docker login -u token -p CRgAAAAA17-2TAyanR1-isuZuMXThkOlueMSXtiE cr.selcloud.ru/dimastephen
	docker push cr.selcloud.ru/dimastephen/chatserver:v0.1

LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3

lint:
	/Users/dmitry/GolandProjects/chatServer/bin/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@latest
get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	mkdir -p pkg/swagger
	make generate-chatServer-api
	$(LOCAL_BIN)/statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png'

generate-chatServer-api:
	mkdir -p pkg/chatServerV1
	protoc --proto_path api/chatServerV1 --proto_path vendor.protogen \
	--go_out=pkg/chatServerV1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chatServerV1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/chatServerV1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	--grpc-gateway_out=pkg/chatServerV1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
	--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
	api/chatServerV1/chatServerV1.proto

include local.env

LOCAL_BIN=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=${MIGRATION_DIR}
LOCAL_MIGRATION_DSN="host=localhost port=${PG_PORT} dbname=${POSTGRES_DB} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD}"


local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v


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



