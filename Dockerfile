FROM golang:1.13 AS builder_deps

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM builder_deps AS builder
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o webhook cmd/main.go

FROM alpine:3.10

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/webhook /app/webhook

ENTRYPOINT ["/app/webhook"]
