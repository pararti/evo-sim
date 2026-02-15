package brain

import "testing"

func TestNetwork_Crossover(t *testing.T) {
	// Two networks with known distinct weights
	n1 := NewNetwork(3, 2, 1)
	n2 := NewNetwork(3, 2, 1)

	// Set all weights to distinguishable values
	for i := range n1.Weights1 {
		for j := range n1.Weights1[i] {
			n1.Weights1[i][j] = -1.0
			n2.Weights1[i][j] = 1.0
		}
	}
	for i := range n1.Weights2 {
		for j := range n1.Weights2[i] {
			n1.Weights2[i][j] = -1.0
			n2.Weights2[i][j] = 1.0
		}
	}

	sawParent1 := false
	sawParent2 := false
	for iter := 0; iter < 100; iter++ {
		child := n1.Crossover(n2)

		// Verify dimensions match
		if len(child.Weights1) != len(n1.Weights1) {
			t.Fatalf("Weights1 row count mismatch")
		}

		for i := range child.Weights1 {
			for j := range child.Weights1[i] {
				w := child.Weights1[i][j]
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
	}

	if !sawParent1 || !sawParent2 {
		t.Errorf("Crossover not mixing weights: p1=%v p2=%v", sawParent1, sawParent2)
	}
}
