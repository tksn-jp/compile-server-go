package main

import (
	"log"
	"net/http"

	"github.com/tksn-jp/compile-server-go/auth"
	"github.com/tksn-jp/compile-server-go/docker"
)

func main() {

	log.Println("Building docker images...")
	imageNum := docker.PrepareImage()

	log.Printf("Image: %d\n", imageNum)

	http.HandleFunc("/auth", auth.Login)
	http.HandleFunc("/exec", Server)
	http.HandleFunc("/ping", Ping)

	log.Println("server ready")

	if err := http.ListenAndServe(":8887", nil); err != nil {
		log.Fatal(err)
	}
}
