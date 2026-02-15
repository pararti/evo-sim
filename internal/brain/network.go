package brain

import (
	"math"
	"math/rand/v2"
)

type Network struct {
	InputSize  int
	HiddenSize int
	OutputSize int
	// Flattened weights for better cache locality
	// weights1 connects input -> hidden (size: InputSize * HiddenSize)
	weights1 []float64
	// weights2 connects hidden -> output (size: HiddenSize * OutputSize)
	weights2 []float64

	// Pre-allocated buffers to reduce GC pressure
	hiddenBuffer []float64
	outputBuffer []float64
}

func NewNetwork(input, hidden, output int) *Network {
	nn := &Network{
		InputSize:    input,
		HiddenSize:   hidden,
		OutputSize:   output,
		weights1:     initWeights(input * hidden),
		weights2:     initWeights(hidden * output),
		hiddenBuffer: make([]float64, hidden),
		outputBuffer: make([]float64, output),
	}
	return nn
}

func (nn *Network) FeedForward(inputs []float64) []float64 {
	// input -> hidden
	// w1 is laid out as [input_0_to_all_hidden, input_1_to_all_hidden, ...]
	// But standard matmul usually iterates [hidden_node][all_inputs].
	// Let's stick to the previous logic:
	// Previous:
	// for i := 0; i < nn.HiddenSize; i++ {
	//   sum := 0.0
	//   for j := 0; j < nn.InputSize; j++ {
	//     sum += inputs[j] * nn.Weights1[j][i]
	//   }
	// }
	//
	// To match this flattened: index = j * HiddenSize + i
	
	for i := 0; i < nn.HiddenSize; i++ {
		sum := 0.0
		for j := 0; j < nn.InputSize; j++ {
			sum += inputs[j] * nn.weights1[j*nn.HiddenSize+i]
		}
		nn.hiddenBuffer[i] = math.Tanh(sum)
	}

	// hidden -> output
	// Previous:
	// for i := 0; i < nn.OutputSize; i++ {
	//   sum := 0.0
	//   for j := 0; j < nn.HiddenSize; j++ {
	//     sum += hiddenOps[j] * nn.Weights2[j][i]
	//   }
	// }
	// Flattened index = j * OutputSize + i

	for i := 0; i < nn.OutputSize; i++ {
		sum := 0.0
		for j := 0; j < nn.HiddenSize; j++ {
			sum += nn.hiddenBuffer[j] * nn.weights2[j*nn.OutputSize+i]
		}
		nn.outputBuffer[i] = math.Tanh(sum)
	}
	
	// Return a copy or the slice? 
	// The original returned a new slice. To be safe and avoid side effects if the caller modifies it,
	// we should return a copy or ensure the caller doesn't hold onto it.
	// Looking at creature.go, it uses values immediately. 
	// But let's return a new slice to match original signature behavior exactly to be safe, 
	// OR (better optimization) rely on the fact that we can copy it out.
	// For max optimization, we return the buffer slice, but we must be careful.
	// creature.go: `dx := output[0]...`. It reads immediately. Safe to return buffer slice.
	return nn.outputBuffer
}

func (nn *Network) Clone() *Network {
	newNet := &Network{
		InputSize:    nn.InputSize,
		HiddenSize:   nn.HiddenSize,
		OutputSize:   nn.OutputSize,
		weights1:     make([]float64, len(nn.weights1)),
		weights2:     make([]float64, len(nn.weights2)),
		hiddenBuffer: make([]float64, nn.HiddenSize),
		outputBuffer: make([]float64, nn.OutputSize),
	}

	copy(newNet.weights1, nn.weights1)
	copy(newNet.weights2, nn.weights2)

	return newNet
}

// Crossover creates a child network by picking each weight from a random parent (uniform crossover).
func (nn *Network) Crossover(other *Network) *Network {
	child := &Network{
		InputSize:    nn.InputSize,
		HiddenSize:   nn.HiddenSize,
		OutputSize:   nn.OutputSize,
		hiddenBuffer: make([]float64, nn.HiddenSize),
		outputBuffer: make([]float64, nn.OutputSize),
	}
	child.weights1 = crossoverSlice(nn.weights1, other.weights1)
	child.weights2 = crossoverSlice(nn.weights2, other.weights2)
	return child
}

func crossoverSlice(a, b []float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		if rand.Float64() < 0.5 {
			result[i] = a[i]
		} else {
			result[i] = b[i]
		}
	}
	return result
}

func (nn *Network) Mutate(rate, strength float64) {
	mutateSlice(nn.weights1, rate, strength)
	mutateSlice(nn.weights2, rate, strength)
}

func initWeights(size int) []float64 {
	w := make([]float64, size)
	for i := range w {
		w[i] = rand.Float64()*2.0 - 1.0
	}
	return w
}

func mutateSlice(weights []float64, rate, strength float64) {
	for i := range weights {
		if rand.Float64() < rate {
			change := rand.NormFloat64() * strength
			weights[i] += change

			if weights[i] > 5.0 {
				weights[i] = 5.0
			}
			if weights[i] < -5.0 {
				weights[i] = -5.0
			}
		}
	}
}
