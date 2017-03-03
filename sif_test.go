package main

import (
	"fmt"
	"testing"
)

func TestScanFileLineNumbers(t *testing.T) {
	var scanFileTests = []struct {
		filename    string
		pattern     string
		lineNumbers []int
	}{
		{"python.txt", "better", []int{3, 4, 5, 6, 7, 8, 17, 18}},
		{"golang.txt", "interface", []int{5, 7}},
		{"subdir/hello.go", "world", []int{6}},
	}

	for _, tc := range scanFileTests {
		s := New(tc.pattern)
		filename := fmt.Sprintf("./_tests/%s", tc.filename)
		fm, err := s.ScanFile(filename)
		if err != nil {
			t.Fatalf("Error on scan file %s: %s", tc.filename, err)
		}
		for i, match := range fm.Matches {
			expected := tc.lineNumbers[i]
			actual := match.Line
			if expected != actual {
				t.Errorf("expected %d, got %d", expected, actual)
			}
		}
	}
}

func TestScanFileText(t *testing.T) {
	s := New("import")
	expected := s.pattern.ReplaceAllStringFunc(`import "fmt"`, bgYellow)

	fm, err := s.ScanFile("./_tests/subdir/hello.go")
	if err != nil {
		t.Fatalf("Error on scan file %s", err)
	}

	if len(fm.Matches) != 1 {
		t.Errorf("expected 1, got %d", len(fm.Matches))
	}

	if fm.Matches[0].Text != expected {
		t.Errorf("expected %s, got %s", expected, fm.Matches[0].Text)
	}
}

func TestScanDir(t *testing.T) {
	expectedFiles := []string{"golang.txt", "subdir/hello.go"}
	s := New("fmt")

	files, err := s.ScanDir("./_tests")
	if err != nil {
		t.Fatalf("Error on scan dir %s", err)
	}

	if len(files) != 2 {
		t.Errorf("expected 2, got %d", len(files))
	}

	for i, f := range files {
		expected := fmt.Sprintf("_tests/%s", expectedFiles[i])
		if f.Name != expected {
			t.Errorf("expected %s, got %s", expected, f.Name)
		}
	}
}
