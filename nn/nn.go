package nn

import (
	"fmt"
	"io"
	"os"

	"gonum.org/v1/gonum/mat"
)

// NeuralNetwork is struct of artificial neural network
type NeuralNetwork struct {
	inputs        int
	hiddens       int
	outputs       int
	hiddenWeights *mat.Dense
	outputWeights *mat.Dense
	learningRate  float64
}

// Predict give an answer
func (net NeuralNetwork) Predict(inputData []float64) mat.Matrix {
	// forward propagation
	inputs := mat.NewDense(len(inputData), 1, inputData)

	hiddenInputs := dot(net.hiddenWeights, inputs)
	// toConsole("hiddenInputs = ", hiddenInputs)

	hiddenOutputs := apply(sigmoid, hiddenInputs)
	// toConsole("hiddenOutputs = ", hiddenOutputs)

	finalInputs := dot(net.outputWeights, hiddenOutputs)
	// toConsole("finalInputs = ", finalInputs)

	finalOutputs := apply(sigmoid, finalInputs)
	// toConsole("finalOutputs = ", finalOutputs)

	return finalOutputs
}

// Train NeuralNetwork
func (net *NeuralNetwork) Train(inputData []float64, targetData []float64) {
	// forward propagation
	inputs := mat.NewDense(len(inputData), 1, inputData)
	hiddenInputs := dot(net.hiddenWeights, inputs)
	hiddenOutputs := apply(sigmoid, hiddenInputs)
	finalInputs := dot(net.outputWeights, hiddenOutputs)
	finalOutputs := apply(sigmoid, finalInputs)

	// find errors
	targets := mat.NewDense(len(targetData), 1, targetData)
	outputErrors := subtract(targets, finalOutputs)
	hiddenErrors := dot(net.outputWeights.T(), outputErrors)

	// backpropagate
	net.outputWeights = add(net.outputWeights,
		scale(net.learningRate,
			dot(multiply(outputErrors, sigmoidPrime(finalOutputs)),
				hiddenOutputs.T()))).(*mat.Dense)

	net.hiddenWeights = add(net.hiddenWeights,
		scale(net.learningRate,
			dot(multiply(hiddenErrors, sigmoidPrime(hiddenOutputs)),
				inputs.T()))).(*mat.Dense)
}

// Save NeuralNetwork
func (net *NeuralNetwork) Save() {
	dir, _ := currentPath()

	h, err := os.Create(fmt.Sprintf("%v/nn/data/hweights.model", dir))
	defer h.Close()
	if err == nil {
		net.hiddenWeights.MarshalBinaryTo(h)
	}
	o, err := os.Create(fmt.Sprintf("%v/nn/data/oweights.model", dir))
	defer o.Close()
	if err == nil {
		net.outputWeights.MarshalBinaryTo(o)
	}
}

// Load NeuralNetwork
func (net *NeuralNetwork) Load() {
	dir, _ := currentPath()

	h, err := os.Open(fmt.Sprintf("%v/nn/saves/hweights.model", dir))
	defer h.Close()
	if err == nil {
		net.hiddenWeights.Reset()
		net.hiddenWeights.UnmarshalBinaryFrom(h)
	}
	o, err := os.Open(fmt.Sprintf("%v/nn/saves/oweights.model", dir))
	defer o.Close()
	if err == nil {
		net.outputWeights.Reset()
		net.outputWeights.UnmarshalBinaryFrom(o)
	}
	return
}

// ImagePredict get predict from image file
func (net *NeuralNetwork) ImagePredict(file io.Reader) (int, float64, []float64) {
	input := dataFromFile(file)
	output := net.Predict(input)
	col := mat.Col(nil, 0, output)

	answer, accuracy := max(col)

	return answer, accuracy, col
}

// CreateNetwork factory with random weights
func CreateNetwork(input, hidden, output int, rate float64) (net NeuralNetwork) {
	net = NeuralNetwork{
		inputs:       input,
		hiddens:      hidden,
		outputs:      output,
		learningRate: rate,
	}

	net.hiddenWeights = mat.NewDense(net.hiddens, net.inputs, randomArray(net.inputs*net.hiddens, float64(net.inputs)))
	net.outputWeights = mat.NewDense(net.outputs, net.hiddens, randomArray(net.hiddens*net.outputs, float64(net.hiddens)))

	return
}

// MnistTrain train network with mnist data
func MnistTrain(net NeuralNetwork) {
	mnistTrain(&net)
	net.Save()
}

// MnistPredict train network with mnist data
func MnistPredict(net NeuralNetwork) {
	net.Load()
	mnistPredict(&net)
}
