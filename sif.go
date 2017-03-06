package sif

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/sync/errgroup"
)

type Match struct {
	Line int
	Text string
}

type FileMatched struct {
	Name    string
	Matches []Match
}

type Options struct {
	CaseInsensitive bool
}

type Sif struct {
	pattern *regexp.Regexp
	options Options
}

var sem = make(chan int, 10)

func (s *Sif) Scan(path string) ([]*FileMatched, error) {
	f, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if f.IsDir() {
		return s.ScanDir(path)
	}

	fm, err := s.ScanFile(path)
	if fm != nil {
		return []*FileMatched{fm}, err
	}
	return nil, err
}

func (s *Sif) ScanDir(dir string) ([]*FileMatched, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var g errgroup.Group
	ch := make(chan *FileMatched)
	for _, f := range files {
		path := filepath.Join(dir, f.Name())
		if !f.IsDir() {
			g.Go(func() error {
				fm, err := s.ScanFile(path)
				if err != nil {
					return err
				}
				if fm != nil {
					ch <- fm
				}
				return nil
			})
		} else if !ignoreDirs.MatchString(f.Name()) {
			g.Go(func() error {
				fs, err := s.ScanDir(path)
				if err != nil {
					return err
				}
				for _, fm := range fs {
					ch <- fm
				}
				return nil
			})
		}
	}

	go func() {
		g.Wait()
		close(ch)
	}()

	filesMatched := make([]*FileMatched, 0)
	for fm := range ch {
		filesMatched = append(filesMatched, fm)
	}

	if err = g.Wait(); err != nil {
		return nil, err
	}

	return filesMatched, nil
}

func (s *Sif) ScanFile(path string) (*FileMatched, error) {
	sem <- 1
	defer func() { <-sem }()

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
		if s.pattern.MatchString(text) {
			text := s.pattern.ReplaceAllStringFunc(text, bgYellow)
			matches = append(matches, Match{line, text})
		}
		line++
	}

	if len(matches) > 0 {
		return &FileMatched{path, matches}, nil
	}

	return nil, nil
}

func New(pattern string, options Options) *Sif {
	if options.CaseInsensitive {
		pattern = fmt.Sprintf("(?i)%s", pattern)
	}
	p := regexp.MustCompile(pattern)
	return &Sif{p, options}
}
