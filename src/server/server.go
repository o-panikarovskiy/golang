package server

import (
	"log"
	"net/http"

	"../nn"
	"./controllers"
)

// Start http server
func Start(port string, net *nn.NeuralNetwork) {
	http.Handle("/", http.FileServer(http.Dir("src/server/www")))
	http.HandleFunc("/api/predict", controllers.GivePredict(net))

	log.Printf("Start listening at %v port\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
