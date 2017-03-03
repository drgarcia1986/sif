# SIF
[![Build Status](https://travis-ci.org/drgarcia1986/sif.svg)](https://travis-ci.org/drgarcia1986/sif)


**S**earch **I**n **F**iles  
An experimental [ack](https://github.com/petdance/ack2) writed in Go.  

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
