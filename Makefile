service-launch:
	cd chat && make install-deps && make vendor-proto && make generate



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