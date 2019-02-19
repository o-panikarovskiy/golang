package main

import (
	"log"
	"net/http"

	"./controllers"
	"./nn"
)

func main() {
	var net = nn.CreateNetwork(784, 200, 10, 0.1)
	net.Load()
	log.Println("NN loaded")

	http.Handle("/", http.FileServer(http.Dir("www")))
	http.HandleFunc("/api/predict", controllers.GivePredict(&net))

	log.Println("Start listening at 8080 port")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
