package main

import (
	"fmt"
	"math"
)

//Num -
const Num int = 859433

func isPrime(n int) bool {
	if n%2 == 0 || n == 1 {
		return false
	}
	for i := 3; i < int(math.Sqrt(float64(n)))+1; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var n, i = 1, 3
	if isPrime(Num) {
		n++
		for i < Num {
			if isPrime(i) {
				n++
			}

			i += 2
		}
		fmt.Printf("%d is prime number(No.%d)\n", Num, n)
	} else {
		fmt.Printf("%d is not prime number\n", Num)
	}
}
