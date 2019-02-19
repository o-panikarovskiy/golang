package nn

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/gonum/mat"
)

func mnistTrain(net *NeuralNetwork) {
	rand.Seed(time.Now().UTC().UnixNano())
	t1 := time.Now()

	dir, _ := currentPath()
	db := fmt.Sprintf("%v/nn/data/mnist_train.csv", dir)

	for epochs := 0; epochs < 5; epochs++ {
		testFile, openError := os.Open(db)

		if openError != nil {
			panic(openError)
		}

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
		testFile.Close()
	}
	elapsed := time.Since(t1)
	fmt.Printf("\nTime taken to train: %s\n", elapsed)
}

func mnistPredict(net *NeuralNetwork) {
	t1 := time.Now()
	dir, _ := currentPath()
	db := fmt.Sprintf("%v/nn/data/mnist_test.csv", dir)

	checkFile, _ := os.Open(db)
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
