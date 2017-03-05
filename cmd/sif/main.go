package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/drgarcia1986/sif"
	"github.com/fatih/color"
)

var opt sif.Options

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [OPTION]... PATTERN [FILES OR DIRECTORIES]\n\n", os.Args[0])
		fmt.Fprint(os.Stdout, "Search for PATTERN in each source file in the tree from the current\n")
		fmt.Fprint(os.Stdout, "directory on down.  If any files or directories are specified, then\n")
		fmt.Fprint(os.Stdout, "only those files and directories are checked.\n\n")
		fmt.Fprintf(os.Stdout, "Example: %s -i select\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.BoolVar(&opt.CaseInsensitive, "i", false, "Ignore case distinctions in PATTERN")
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
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

func scan(pattern string, dirs ...string) ([]*sif.FileMatched, error) {
	s := sif.New(pattern, opt)
	files := make([]*sif.FileMatched, 0)
	for _, dir := range dirs {
		fs, err := s.Scan(dir)
		if err != nil {
			return nil, err
		}

		if fs != nil {
			files = append(files, fs...)
		}
	}
	return files, nil
}

func show(files ...*sif.FileMatched) {
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
