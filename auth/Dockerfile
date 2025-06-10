FROM golang:1.24.2-alpine AS builder

COPY . /auth/
WORKDIR /auth/

RUN go mod download
RUN go build --platform linux/amd64 -o /bin/server cmd/main.go

FROM alpine:latest
COPY --from=builder /bin/server .

CMD ["./server"]