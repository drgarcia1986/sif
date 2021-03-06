# SIF
[![Build Status](https://travis-ci.org/drgarcia1986/sif.svg)](https://travis-ci.org/drgarcia1986/sif)
[![Go Report Card](https://goreportcard.com/badge/drgarcia1986/sif)](https://goreportcard.com/report/drgarcia1986/sif)
[![codecov](https://codecov.io/gh/drgarcia1986/sif/branch/master/graph/badge.svg)](https://codecov.io/gh/drgarcia1986/sif)


**S**earch **I**n **F**iles  
An experimental [ack](https://github.com/petdance/ack2) written in Go.  

## Installation
If you have a Golang environment setup, you can simply run:
```
$ go get -u github.com/drgarcia1986/sif/cmd/sif
```
Or get binaries on [releases](https://github.com/drgarcia1986/sif/releases)

> For windows users: Colors works only in PowerShell

## Example
Run against repo dir:
```
$ sif better
_tests/golang.txt
9: A little copying is better than a little dependency.
14: Clear is better than clever.

_tests/python.txt
3: Beautiful is better than ugly.
4: Explicit is better than implicit.
5: Simple is better than complex.
6: Complex is better than complicated.
7: Flat is better than nested.
8: Sparse is better than dense.
17: Now is better than never.
18: Although never is often better than *right* now.

sif_test.go
14:             {"python.txt", "better", []int{3, 4, 5, 6, 7, 8, 17, 18}},
```
Same search with `ack`, `grep` and `sif`, time comparison:
```
$ time ack better
...
        0.11 real         0.07 user         0.01 sys
```
```
$ time grep better -rn *
...
        0.04 real         0.03 user         0.00 sys
```
```
$ time sif better
...
        0.01 real         0.00 user         0.00 sys
```

## Library Usage Example
```go
package main

import (
	"fmt"

	"github.com/drgarcia1986/sif"
)

func main() {
	s, _ := sif.New("fmt", sif.Options{CaseInsensitive: false})
	fm, err := s.ScanFile("./main.go")
	if err != nil {
		panic(err)
	}
	if fm != nil {
		for _, m := range fm.Matches {
			fmt.Printf("Line: %d, Text: %s\n", m.Line, m.Text)
		}
	}
}
```
