package nn

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/gonum/mat"
)

func runTestEpoh(net *NeuralNetwork, path string) {
	testFile, openError := os.Open(path)

	if openError != nil {
		log.Fatal(openError)
	}

	defer testFile.Close()

	r := csv.NewReader(bufio.NewReader(testFile))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		inputs := make([]float64, net.inputs)
		for i := range inputs {
			x, _ := strconv.ParseFloat(record[i], 64)
			inputs[i] = (x / 255.0 * 0.99) + 0.01
		}

		targets := make([]float64, 10)
		for i := range targets {
			targets[i] = 0.01
		}
		x, _ := strconv.Atoi(record[0])
		targets[x] = 0.99

		net.Train(inputs, targets)
	}
}

func mnistTrain(net *NeuralNetwork, path string) {
	rand.Seed(time.Now().UTC().UnixNano())
	t1 := time.Now()

	for epochs := 0; epochs < 5; epochs++ {
		runTestEpoh(net, path)
	}

	elapsed := time.Since(t1)
	fmt.Printf("\nTime taken to train: %s\n", elapsed)
}

func mnistPredict(net *NeuralNetwork, path string) {
	t1 := time.Now()
	checkFile, openError := os.Open(path)

	if openError != nil {
		log.Fatal(openError)
	}

	defer checkFile.Close()

	score := 0
	r := csv.NewReader(bufio.NewReader(checkFile))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		inputs := make([]float64, net.inputs)

		for i := range inputs {
			if i == 0 {
				inputs[i] = 1.0
			}
			x, _ := strconv.ParseFloat(record[i], 64)
			inputs[i] = (x / 255.0 * 0.99) + 0.01
		}

		outputs := net.Predict(inputs)
		best, _ := max(mat.Col(nil, 0, outputs))
		target, _ := strconv.Atoi(record[0])

		if best == target {
			score++
		}
	}

	elapsed := time.Since(t1)
	fmt.Printf("Time taken to check: %s\n", elapsed)
	fmt.Println("score:", score)
}
