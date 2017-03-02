package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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

var (
	colorBgYellow = color.New(color.BgYellow, color.FgBlack).SprintFunc()
	bgYellow      = func(s string) string { return colorBgYellow(s) }
	yellow        = color.New(color.Bold, color.FgYellow).SprintFunc()
	green         = color.New(color.Bold, color.FgGreen)
	ignoreDirs    = regexp.MustCompile(`\.(svn|git*)`)
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("usage: sif <pattern> [directories...]")
		os.Exit(1)
	}

	pattern := regexp.MustCompile(fmt.Sprintf("(?i)%s", args[0]))

	dirs := make([]string, 0)
	if len(args) < 2 {
		dirs = append(dirs, "./")
	} else {
		dirs = append(dirs, args[1:]...)
	}

	files := make([]*FileMatched, 0)
	for _, dir := range dirs {
		fs, err := scanDir(dir, pattern)
		if err != nil {
			fmt.Printf("error to search for pattern: %s", err)
			os.Exit(1)
		}
		files = append(files, fs...)
	}

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

func scanDir(dir string, pattern *regexp.Regexp) ([]*FileMatched, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filesMatched := make([]*FileMatched, 0)
	for _, f := range files {
		path := filepath.Join(dir, f.Name())
		if !f.IsDir() {
			fm, err := scanFile(path, pattern)
			if err != nil {
				return nil, err
			}
			if fm != nil {
				filesMatched = append(filesMatched, fm)
			}
		} else if !ignoreDirs.MatchString(f.Name()) {
			fs, err := scanDir(path, pattern)
			if err != nil {
				return nil, err
			}
			filesMatched = append(filesMatched, fs...)
		}
	}

	return filesMatched, nil
}

func scanFile(path string, pattern *regexp.Regexp) (*FileMatched, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if b, err := isBinary(file); b || err != nil {
		return nil, err
	}

	line := 1
	matches := make([]Match, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if pattern.MatchString(text) {
			text := pattern.ReplaceAllStringFunc(text, bgYellow)
			matches = append(matches, Match{line, text})
		}
		line++
	}

	if len(matches) > 0 {
		return &FileMatched{path, matches}, nil
	}

	return nil, nil
}

func isBinary(file *os.File) (bool, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return true, err
	}
	file.Seek(0, 0)
	return http.DetectContentType(buffer[:n]) == "application/octet-stream", nil
}
