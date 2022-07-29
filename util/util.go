package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// ReadLines 一行一行的读取文件，追加到一个字符串数组中
func ReadLines(file string) ([]string, error) {
	return ReadSplitter(file, '\n')
}

// ReadSplitter 使用指定的分隔符读取文件
func ReadSplitter(file string, splitter byte) (lines []string, err error) {
	fin, err := os.Open(file)
	if err != nil {
		return
	}

	r := bufio.NewReader(fin)
	for {
		line, err := r.ReadString(splitter)
		if err == io.EOF {
			break
		}
		line = strings.Replace(line, string(splitter), "", -1)
		lines = append(lines, line)
	}
	return
}
