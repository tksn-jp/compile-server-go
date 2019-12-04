FROM golang:1.13.4-alpine

WORKDIR /src

CMD go build -ldflags='-w -s' -o app . && ./app
