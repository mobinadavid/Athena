FROM golang:alpine AS builder
LABEL maintainer="Ali Khedmati"
# ENVs
ENV CGO_ENABLED 0 \
    GOPATH /go \
    GOCACHE /go-build \
    GOOS linux \
    GOARCH=amd64 \
    GO111MODULE=on
# Copy
WORKDIR /athena
COPY . /athena
# Install Dependencies
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    go mod download

# Build the binary
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go build -o athena main.go

CMD ["./athena", "app", "bootstrap"]
