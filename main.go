package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/fatih/color"
)

type Match struct {
	Line int
	Text string
}

type FileMatched struct {
	Name    string
	Matches []Match
}

var colorBgYellow = color.New(color.BgYellow, color.FgBlack).SprintFunc()
var bgYellow = func(s string) string { return colorBgYellow(s) }
var yellow = color.New(color.Bold, color.FgYellow).SprintFunc()
var green = color.New(color.Bold, color.FgGreen)

func main() {
	pattern := regexp.MustCompile("foo")

	files, err := scanDir("./", pattern)
	if err != nil {
		fmt.Printf("error to search for pattern: %s", err)
		os.Exit(1)
	}

	for _, f := range files {
		green.Println(f.Name)
		for _, match := range f.Matches {
			fmt.Printf("%s: %s\n", yellow(match.Line), match.Text)
		}
	}
}

func scanDir(dir string, pattern *regexp.Regexp) ([]*FileMatched, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filesMatched := make([]*FileMatched, 0)
	for _, f := range files {
		fm, err := scanFile(f.Name(), pattern)
		if err != nil {
			return nil, err
		}
		if fm != nil {
			filesMatched = append(filesMatched, fm)
		}
	}

	return filesMatched, nil
}

func scanFile(filename string, pattern *regexp.Regexp) (*FileMatched, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	matches := make([]Match, 0)
	lineCount := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if pattern.MatchString(text) {
			text := pattern.ReplaceAllStringFunc(text, bgYellow)
			matches = append(matches, Match{lineCount, text})
		}
		lineCount++
	}

	if len(matches) > 0 {
		return &FileMatched{filename, matches}, nil
	}

	return nil, nil
}
