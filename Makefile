service-launch:
	cd chat && make install-deps && make vendor-proto && make generate
	cd auth && make install-deps && make vendor-proto && make generate
	docker compose up -d


grpc-load-chat-test:
	ghz \
		--proto chat/api/chatServerV1/chatServerV1.proto \
		--import-paths chat/vendor.protogen \
		--call chatServerV1.ChatServer/Create \
		--data '{"usernames":1}' \
		--rps 50 \
		--total 200 \
		--insecure \
		localhost:50051

grpc-load-auth-test:
	ghz \
		--proto auth/api/authV1/auth.proto \
		--import-paths auth/vendor.protogen \
		--call authV1.Auth/GetRefreshToken \
		--data '{"refreshToken":"fnsdvfjsiofvps"}' \
		--rps 1000 \
		--total 5000 \
		--insecure \
		localhost:50052