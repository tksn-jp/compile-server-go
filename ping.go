package main

import "net/http"

func Ping(res http.ResponseWriter, req *http.Request) {
	bytes := []byte("pongpong")
	_, _ = res.Write(bytes)
}
