package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("usage: sif <pattern> [directories...]")
		os.Exit(1)
	}

	dirs := make([]string, 0)
	if len(args) < 2 {
		dirs = append(dirs, "./")
	} else {
		dirs = append(dirs, args[1:]...)
	}

	files, err := scan(args[0], dirs...)
	if err != nil {
		fmt.Printf("error to search for pattern: %s", err)
		os.Exit(1)
	}

	show(files...)
}

func scan(pattern string, dirs ...string) ([]*FileMatched, error) {
	s := New(pattern)
	files := make([]*FileMatched, 0)
	for _, dir := range dirs {
		f, err := os.Stat(dir)
		if err != nil {
			return nil, err
		}

		if f.IsDir() {
			fs, err := s.ScanDir(dir)
			if err != nil {
				return nil, err
			}
			files = append(files, fs...)
		} else {
			fm, err := s.ScanFile(dir)
			if err != nil {
				return nil, err
			}
			if fm != nil {
				files = append(files, fm)
			}
		}
	}
	return files, nil
}

func show(files ...*FileMatched) {
	green := color.New(color.Bold, color.FgGreen)
	yellow := color.New(color.Bold, color.FgYellow).SprintFunc()
	for i, f := range files {
		green.Println(f.Name)
		for _, match := range f.Matches {
			fmt.Printf("%s: %s\n", yellow(match.Line), match.Text)
		}
		if i+1 < len(files) {
			fmt.Println()
		}
	}
}
