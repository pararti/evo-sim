package brain

import "testing"

func TestNetwork_Crossover(t *testing.T) {
	// Two networks with known distinct weights
	n1 := NewNetwork(3, 2, 1)
	n2 := NewNetwork(3, 2, 1)

	// Set all weights to distinguishable values
	for i := range n1.weights1 {
		n1.weights1[i] = -1.0
		n2.weights1[i] = 1.0
	}
	for i := range n1.weights2 {
		n1.weights2[i] = -1.0
		n2.weights2[i] = 1.0
	}

	sawParent1 := false
	sawParent2 := false
	for iter := 0; iter < 100; iter++ {
		child := n1.Crossover(n2)

		// Verify dimensions match
		if len(child.weights1) != len(n1.weights1) {
			t.Fatalf("weights1 count mismatch")
		}

		for i := range child.weights1 {
			w := child.weights1[i]
			if w != -1.0 && w != 1.0 {
				t.Fatalf("Child weight %f is neither parent", w)
			}
			if w == -1.0 {
				sawParent1 = true
			}
			if w == 1.0 {
				sawParent2 = true
			}
		}
	}

	if !sawParent1 || !sawParent2 {
		t.Errorf("Crossover not mixing weights: p1=%v p2=%v", sawParent1, sawParent2)
	}
}

func TestNetwork_ElmanRecurrence(t *testing.T) {
	nn := NewNetwork(2, 3, 1)

	// Zero all weights for predictability
	for i := range nn.weights1 {
		nn.weights1[i] = 0
	}
	for i := range nn.weights2 {
		nn.weights2[i] = 0
	}

	// Set a recurrent weight: hidden state[0] -> hidden[0]
	// Recurrent row starts at InputSize * HiddenSize
	// Index: (InputSize + 0) * HiddenSize + 0 = 2*3 + 0 = 6
	nn.weights1[6] = 1.0
	// Set weights2[0] = 1.0 so output reflects hidden[0]
	nn.weights2[0] = 1.0
	// Set a sensory weight: input[0] -> hidden[0]
	nn.weights1[0] = 1.0

	// First call: no recurrent state, input[0]=0.5
	raw1 := nn.FeedForward([]float64{0.5, 0.0})
	val1 := raw1[0] // Capture value before next call overwrites buffer

	// Second call: same input, but now hidden state is non-zero
	raw2 := nn.FeedForward([]float64{0.5, 0.0})
	val2 := raw2[0]

	// Output should differ because of recurrent state
	if val1 == val2 {
		t.Errorf("Elman recurrence should produce different outputs for same input on consecutive ticks: %f vs %f", val1, val2)
	}

	// Second output should be larger in magnitude (more activation from recurrence)
	if val2 <= val1 {
		t.Errorf("Expected second output (%f) > first output (%f) due to positive recurrence", val2, val1)
	}
}

func TestNetwork_CloneWithResize(t *testing.T) {
	nn := NewNetwork(3, 4, 2)

	// Shrink
	smaller := nn.CloneWithResize(2)
	if smaller.HiddenSize != 2 {
		t.Errorf("Expected hidden size 2, got %d", smaller.HiddenSize)
	}
	if len(smaller.weights1) != (3+2)*2 {
		t.Errorf("Expected weights1 size %d, got %d", (3+2)*2, len(smaller.weights1))
	}

	// Grow
	larger := nn.CloneWithResize(6)
	if larger.HiddenSize != 6 {
		t.Errorf("Expected hidden size 6, got %d", larger.HiddenSize)
	}

	// Should still produce valid output
	out := smaller.FeedForward([]float64{0.1, 0.2, 0.3})
	if len(out) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(out))
	}
}
