package main

import (
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/fatih/color"
)

var (
	colorBgYellow = color.New(color.BgYellow, color.FgBlack).SprintFunc()
	bgYellow      = func(s string) string { return colorBgYellow(s) }
	ignoreDirs    = regexp.MustCompile(`\.(svn|git*)`)
)

func isBinary(file *os.File) (bool, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return true, err
	}
	file.Seek(0, 0)
	return http.DetectContentType(buffer[:n]) == "application/octet-stream", nil
}
