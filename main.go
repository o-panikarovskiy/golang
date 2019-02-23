package main

import (
	"flag"
	"fmt"
	"os"

	"./src/nn"
	"./src/server"
)

func loadNetwork(net *nn.NeuralNetwork) {
	err := net.Load()

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Println("NN loaded")
}

func main() {
	net := nn.CreateNetwork(784, 200, 10, 0.1)

	command := flag.String("cmd", "", "train | predict | server")
	port := flag.String("port", "8080", "Port for listening")
	file := flag.String("file", "", "MNIST CSV file for train or predict")
	flag.Parse()

	if *command == "train" {
		nn.MnistTrain(&net, *file)
		net.Save()
	} else if *command == "predict" {
		loadNetwork(&net)
		nn.MnistPredict(&net, *file)
	} else if *command == "server" {
		loadNetwork(&net)
		server.Start(":"+*port, &net)
	}
}
