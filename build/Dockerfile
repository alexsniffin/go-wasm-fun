############# Build the binary #############
FROM golang:1.15-alpine AS builder

RUN apk update && apk upgrade && apk --no-cache add curl

# Install latest gcc
RUN apk add build-base

# Copy local source
COPY . $GOPATH/src/github.com/alexsniffin/website.git/
WORKDIR $GOPATH/src/github.com/alexsniffin/website.git/

# Build binarys
RUN GOOS=linux go build -a -o /app/server ./cmd/server
RUN GOOS=js GOARCH=wasm go build -a -o /app/wasm ./cmd/wasm

COPY ./assets /app/assets/

############# Build the image #############
FROM alpine:3

RUN apk update && apk upgrade

# Copy config and binary
WORKDIR /app/
COPY ./configs/server.yaml ./configs/
COPY --from=builder /app/ .

# Entrypoint and port
CMD ["./server"]
EXPOSE 80