package brain

import (
	"math"
	"math/rand/v2"
)

type Network struct {
	InputSize  int // Sensory input size (not including recurrent context)
	HiddenSize int
	OutputSize int
	// Flattened weights for better cache locality
	// weights1 connects (input + hiddenState) -> hidden (size: (InputSize + HiddenSize) * HiddenSize)
	weights1 []float64
	// weights2 connects hidden -> output (size: HiddenSize * OutputSize)
	weights2 []float64

	// Elman recurrent memory: previous tick's hidden layer output
	hiddenState []float64

	// Pre-allocated buffers to reduce GC pressure
	hiddenBuffer []float64
	outputBuffer []float64
}

// totalInputSize returns the effective input size including recurrent context.
func (nn *Network) totalInputSize() int {
	return nn.InputSize + nn.HiddenSize
}

func NewNetwork(input, hidden, output int) *Network {
	totalInput := input + hidden
	nn := &Network{
		InputSize:    input,
		HiddenSize:   hidden,
		OutputSize:   output,
		weights1:     initWeights(totalInput * hidden),
		weights2:     initWeights(hidden * output),
		hiddenState:  make([]float64, hidden),
		hiddenBuffer: make([]float64, hidden),
		outputBuffer: make([]float64, output),
	}
	return nn
}

func (nn *Network) FeedForward(inputs []float64) []float64 {
	// Elman network: effective input = [sensory_inputs, hiddenState]
	totalIn := nn.totalInputSize()

	// (input + hiddenState) -> hidden
	// Layout: j * HiddenSize + i, where j indexes the total input
	for i := 0; i < nn.HiddenSize; i++ {
		sum := 0.0
		// Sensory inputs
		for j := 0; j < nn.InputSize; j++ {
			sum += inputs[j] * nn.weights1[j*nn.HiddenSize+i]
		}
		// Recurrent context (previous hidden state)
		for j := nn.InputSize; j < totalIn; j++ {
			sum += nn.hiddenState[j-nn.InputSize] * nn.weights1[j*nn.HiddenSize+i]
		}
		nn.hiddenBuffer[i] = math.Tanh(sum)
	}

	// Save hidden output as state for next tick
	copy(nn.hiddenState, nn.hiddenBuffer)

	// hidden -> output
	for i := 0; i < nn.OutputSize; i++ {
		sum := 0.0
		for j := 0; j < nn.HiddenSize; j++ {
			sum += nn.hiddenBuffer[j] * nn.weights2[j*nn.OutputSize+i]
		}
		nn.outputBuffer[i] = math.Tanh(sum)
	}

	return nn.outputBuffer
}

func (nn *Network) Clone() *Network {
	newNet := &Network{
		InputSize:    nn.InputSize,
		HiddenSize:   nn.HiddenSize,
		OutputSize:   nn.OutputSize,
		weights1:     make([]float64, len(nn.weights1)),
		weights2:     make([]float64, len(nn.weights2)),
		hiddenState:  make([]float64, nn.HiddenSize), // Zeroed — newborns have no memory
		hiddenBuffer: make([]float64, nn.HiddenSize),
		outputBuffer: make([]float64, nn.OutputSize),
	}

	copy(newNet.weights1, nn.weights1)
	copy(newNet.weights2, nn.weights2)

	return newNet
}

// Crossover creates a child network by picking each weight from a random parent (uniform crossover).
// Both parents must have the same hidden size.
func (nn *Network) Crossover(other *Network) *Network {
	child := &Network{
		InputSize:    nn.InputSize,
		HiddenSize:   nn.HiddenSize,
		OutputSize:   nn.OutputSize,
		hiddenState:  make([]float64, nn.HiddenSize), // Zeroed
		hiddenBuffer: make([]float64, nn.HiddenSize),
		outputBuffer: make([]float64, nn.OutputSize),
	}
	child.weights1 = crossoverSlice(nn.weights1, other.weights1)
	child.weights2 = crossoverSlice(nn.weights2, other.weights2)
	return child
}

// CloneWithResize creates a copy of the network with a potentially different hidden size.
// Matching weights are copied; extras are randomly initialized.
func (nn *Network) CloneWithResize(newHiddenSize int) *Network {
	newNet := NewNetwork(nn.InputSize, newHiddenSize, nn.OutputSize)

	minH := nn.HiddenSize
	if newHiddenSize < minH {
		minH = newHiddenSize
	}

	// Copy weights1: total input = InputSize + HiddenSize
	// For sensory input rows (j < InputSize): layout j*HiddenSize+i
	for j := 0; j < nn.InputSize; j++ {
		for i := 0; i < minH; i++ {
			newNet.weights1[j*newHiddenSize+i] = nn.weights1[j*nn.HiddenSize+i]
		}
	}
	// For recurrent rows (j = InputSize .. InputSize+minH-1):
	// Old layout: (InputSize + k) * oldHidden + i, where k is the recurrent index
	// New layout: (InputSize + k) * newHidden + i
	for k := 0; k < minH; k++ {
		for i := 0; i < minH; i++ {
			newNet.weights1[(nn.InputSize+k)*newHiddenSize+i] = nn.weights1[(nn.InputSize+k)*nn.HiddenSize+i]
		}
	}

	// Copy weights2
	for j := 0; j < minH; j++ {
		for i := 0; i < nn.OutputSize; i++ {
			newNet.weights2[j*nn.OutputSize+i] = nn.weights2[j*nn.OutputSize+i]
		}
	}

	return newNet
}

// CrossoverWithResize creates a child network from two parents that may have different hidden sizes.
// The child uses the specified hiddenSize. Matching weights are crossed over; extras are randomized.
func (nn *Network) CrossoverWithResize(other *Network, childHiddenSize int) *Network {
	child := NewNetwork(nn.InputSize, childHiddenSize, nn.OutputSize)

	childTotalInput := nn.InputSize + childHiddenSize

	// weights1: for each total input j, for each hidden i
	for j := 0; j < childTotalInput; j++ {
		for i := 0; i < childHiddenSize; i++ {
			// Determine the corresponding index in each parent
			var w1InRange, w2InRange bool
			var v1, v2 float64

			if j < nn.InputSize {
				// Sensory input row — same row index in both parents
				if i < nn.HiddenSize {
					w1InRange = true
					v1 = nn.weights1[j*nn.HiddenSize+i]
				}
				if i < other.HiddenSize {
					w2InRange = true
					v2 = other.weights1[j*other.HiddenSize+i]
				}
			} else {
				// Recurrent row — index k = j - InputSize
				k := j - nn.InputSize
				if k < nn.HiddenSize && i < nn.HiddenSize {
					w1InRange = true
					v1 = nn.weights1[(nn.InputSize+k)*nn.HiddenSize+i]
				}
				if k < other.HiddenSize && i < other.HiddenSize {
					w2InRange = true
					v2 = other.weights1[(other.InputSize+k)*other.HiddenSize+i]
				}
			}

			if w1InRange && w2InRange {
				if rand.Float64() < 0.5 {
					child.weights1[j*childHiddenSize+i] = v1
				} else {
					child.weights1[j*childHiddenSize+i] = v2
				}
			} else if w1InRange {
				child.weights1[j*childHiddenSize+i] = v1
			} else if w2InRange {
				child.weights1[j*childHiddenSize+i] = v2
			}
			// else: keep random init from NewNetwork
		}
	}

	// weights2: for each hidden j, for each output i
	for j := 0; j < childHiddenSize; j++ {
		for i := 0; i < nn.OutputSize; i++ {
			var w1InRange, w2InRange bool
			var v1, v2 float64
			if j < nn.HiddenSize {
				w1InRange = true
				v1 = nn.weights2[j*nn.OutputSize+i]
			}
			if j < other.HiddenSize {
				w2InRange = true
				v2 = other.weights2[j*other.OutputSize+i]
			}
			if w1InRange && w2InRange {
				if rand.Float64() < 0.5 {
					child.weights2[j*nn.OutputSize+i] = v1
				} else {
					child.weights2[j*nn.OutputSize+i] = v2
				}
			} else if w1InRange {
				child.weights2[j*nn.OutputSize+i] = v1
			} else if w2InRange {
				child.weights2[j*nn.OutputSize+i] = v2
			}
		}
	}

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
