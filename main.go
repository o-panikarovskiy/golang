package main

import (
	"flag"
	"log"
	"os"

	"./src/nn"
	"./src/server"
)

func main() {
	var net = nn.CreateNetwork(784, 200, 10, 0.1)
	err := net.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("NN Loaded")

	command := os.Args[1]
	port := flag.String("port", "", "Port for listening")

	flag.Parse()

	if *port == "" {
		*port = ":8080"
	} else {
		*port = ":" + *port
	}

	// train or mass predict to determine the effectiveness of the trained network
	switch command {
	case "train":
		nn.MnistTrain(&net)
		net.Save()
	case "predict":
		net.Load()
		nn.MnistPredict(&net)
	case "server":
		server.Start(*port, &net)
	default:
		// don't do anything
	}
}
