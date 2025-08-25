FROM golang:1.24.2-alpine AS builder
RUN apk add git openssh
RUN mkdir -p /root/.ssh
COPY ./id_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

COPY auth/ auth/
WORKDIR auth/

RUN ls -R

RUN go mod download
RUN go build  -o /bin/server ./cmd/main.go

FROM alpine:latest
COPY --from=builder /bin/server .
COPY ./auth/local.env .
COPY ./.env .
COPY ./auth/roles.json .
CMD ["./server"]