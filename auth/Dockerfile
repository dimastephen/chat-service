FROM golang:1.24.2-alpine AS builder

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