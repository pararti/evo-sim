package brain

import (
	"math"
	"math/rand/v2"
)

type Network struct {
	InputSize  int
	HiddenSize int
	OutputSize int
	//[from][to]
	//connect input and hidden
	Weights1 [][]float64
	Weights2 [][]float64
}

func NewNetwork(input, hidden, output int) *Network {
	nn := &Network{
		InputSize:  input,
		HiddenSize: hidden,
		OutputSize: output,
	}

	nn.Weights1 = initWeights(input, hidden)
	nn.Weights2 = initWeights(hidden, output)
	return nn
}

func (nn *Network) FeedForward(inputs []float64) []float64 {
	//input -> hidden
	hiddenOps := make([]float64, nn.HiddenSize)
	for i := 0; i < nn.HiddenSize; i++ {
		sum := 0.0
		for j := 0; j < nn.InputSize; j++ {
			sum += inputs[j] * nn.Weights1[j][i]
		}
		hiddenOps[i] = math.Tanh(sum)
	}

	//hidden -> output
	finalOps := make([]float64, nn.OutputSize)
	for i := 0; i < nn.OutputSize; i++ {
		sum := 0.0
		for j := 0; j < nn.HiddenSize; j++ {
			sum += hiddenOps[j] * nn.Weights2[j][i]
		}
		finalOps[i] = math.Tanh(sum)
	}
	return finalOps
}

func (nn *Network) Clone() *Network {
	newNet := &Network{
		InputSize:  nn.InputSize,
		HiddenSize: nn.HiddenSize,
		OutputSize: nn.OutputSize,
	}

	newNet.Weights1 = make([][]float64, len(nn.Weights1))
	for i := range nn.Weights1 {
		newNet.Weights1[i] = make([]float64, len(nn.Weights1[i]))
		copy(newNet.Weights1[i], nn.Weights1[i])
	}

	newNet.Weights2 = make([][]float64, len(nn.Weights2))
	for i := range nn.Weights2 {
		newNet.Weights2[i] = make([]float64, len(nn.Weights2[i]))
		copy(newNet.Weights2[i], nn.Weights2[i])
	}

	return newNet
}

func (nn *Network) Mutate(rate, strength float64) {
	mutateMatrix(nn.Weights1, rate, strength)
	mutateMatrix(nn.Weights2, rate, strength)
}

func initWeights(rows, cols int) [][]float64 {
	w := make([][]float64, rows)
	for i := range w {
		w[i] = make([]float64, cols)
		//random weights init
		for j := range w[i] {
			w[i][j] = rand.Float64()
		}
	}
	return w
}

func mutateMatrix(weights [][]float64, rate, strength float64) {
	for i := range weights {
		for j := range weights[i] {
			if rand.Float64() < rate {
				change := rand.NormFloat64() * strength
				weights[i][j] += change

				if weights[i][j] > 5.0 {
					weights[i][j] = 5.0
				}
				if weights[i][j] < -5.0 {
					weights[i][j] = -5.0
				}
			}
		}
	}
}
