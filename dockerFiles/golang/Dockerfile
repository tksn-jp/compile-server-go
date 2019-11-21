FROM golang:latest

WORKDIR /go/src

CMD go build -ldflags='-w -s' -o app . && ./app
