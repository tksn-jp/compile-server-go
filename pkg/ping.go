package pkg

import "net/http"

func Ping(res http.ResponseWriter, req *http.Request) {
	bytes := []byte("pong")
	_, _ = res.Write(bytes)
}
