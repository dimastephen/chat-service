FROM golang:1.23.1-alpine AS builder

COPY . /github.com/dimastephen/chatServer/
WORKDIR /github.com/dimastephen/chatServer/

RUN go mod download
RUN go build -o ./bin/linux_server cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/dimastephen/chatServer/bin/linux_server .

CMD ["./linux_server"]