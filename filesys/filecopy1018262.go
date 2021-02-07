package main

import (
	"fmt"
	"os"
	"syscall"
	"unicode"
)

func main() {
	//Open input file
	fi, err := syscall.Open(os.Args[1], syscall.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := syscall.Close(fi); err != nil {
			panic(err)
		}
	}()

	//Open output file
	fo, err := syscall.Open(os.Args[2], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := syscall.Close(fo); err != nil {
			panic(err)
		}
	}()

	//Get status of input file
	var sti syscall.Stat_t
	if err := syscall.Fstat(fi, &sti); err != nil {
		panic(err)
	}

	//Read input file by syscall.Mmap
	mmap, err := syscall.Mmap(fi, 0, int(sti.Size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	//Parse rune
	dist := make([]byte, 0)
	for _, s := range mmap {
		if (unicode.IsSpace(rune(s)) || unicode.IsLetter(rune(s)) || unicode.IsNumber(rune(s))) && s < 126 {
			dist = append(dist, byte(s))
		}
	}

	//Unmap
	if err := syscall.Munmap(mmap); err != nil {
		panic(err)
	}

	//Write to output file
	if _, err := syscall.Write(fo, dist); err != nil {
		panic(err)
	}

	//Print size of output file
	var sto syscall.Stat_t
	if err := syscall.Fstat(fo, &sto); err != nil {
		panic(err)
	}
	fmt.Println(sto.Size)
}
