package main

import (
	"flag"
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

	port := flag.String("port", "", "Port for listening")
	flag.Parse()

	if *port == "" {
		*port = ":8080"
	} else {
		*port = ":" + *port
	}

	log.Printf("Start listening at %v port\n", *port)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
