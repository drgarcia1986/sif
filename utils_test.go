package sif

import (
	"fmt"
	"os"
	"testing"
)

func TestIsBinary(t *testing.T) {
	var isBinaryTests = []struct {
		filename string
		expected bool
	}{
		{"python.txt", false},
		{"binary_file", true},
	}

	for _, tc := range isBinaryTests {
		path := fmt.Sprintf("./_tests/%s", tc.filename)
		file, err := os.Open(path)
		if err != nil {
			t.Fatalf("Error on open file %s", err)
		}
		defer file.Close()

		actual, err := isBinary(file)
		if err != nil {
			t.Fatalf("Error on check file %s", err)
		}

		if actual != tc.expected {
			t.Errorf("expected %t, got %t", tc.expected, actual)
		}
	}
}
