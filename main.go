package main

import (
	"log"
	"net/http"

	"github.com/tksn-jp/compile-server-go/pkg"
	"github.com/tksn-jp/compile-server-go/pkg/auth"
	"github.com/tksn-jp/compile-server-go/pkg/file"
)

func main() {

	//id, err := docker.BuildContainer("golang", "cl-golang")
	//if err != nil {
	//	if err != errors.New("") {
	//		log.Fatal(err)
	//	}
	//}
	//log.Printf("container builded: %s", id)
	//docker.StartContainer(id)

	log.Println("server ready")

	http.HandleFunc("/auth", auth.Login)
	http.HandleFunc("/compile", file.Compile)
	http.HandleFunc("/ping", pkg.Ping)

	log.Fatal(http.ListenAndServe(":8888", nil))
}
