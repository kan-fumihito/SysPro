package main

import (
	"bufio"
	"io"
	"os"
)

func main() {
	fi, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	br := bufio.NewReader(fi)
	fo, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	bw := bufio.NewWriter(fo)

	buf := make([]byte, 1024)
	for {
		n, err := br.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		if _, err := bw.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
	if err := bw.Flush(); err != nil {
		panic(err)
	}
}
