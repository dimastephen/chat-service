FROM golang:1.24.2-alpine AS builder

RUN apk add git openssh
RUN mkdir -p /root/.ssh
COPY ./id_rsa ./root/.ssh/id_rsa
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

COPY ./chat/ ./chat/
WORKDIR ./chat/

RUN ls -R

RUN go mod download
RUN go build -o /bin/chat_server cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /bin/chat_server .
COPY ./chat/local.env .
COPY ./compose.env .


CMD ["./chat_server"]