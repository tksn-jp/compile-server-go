package main

import (
	"log"
	"net/http"

	"github.com/tksn-jp/compile-server-go/auth"
	"github.com/tksn-jp/compile-server-go/docker"
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
	log.Println("Building docker images...")
	imageNum := docker.PrepareImage()

	log.Printf("Image: %d\n", imageNum)
	log.Println("server ready")

	http.HandleFunc("/auth", auth.Login)
	http.HandleFunc("/compile", Server)
	http.HandleFunc("/ping", Ping)

	log.Fatal(http.ListenAndServe(":8888", nil))
}
