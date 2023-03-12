package main

import (
	"fmt"
	"math/rand"
)

func modexp(g, p, b int) int {
	res := 1
	for i := 0; i < b; i++ {
		res = (res * g) % p
	}
	return res
}

func main() {
	var g, p int
	fmt.Scanf("g is %d and p is %d", &g, &p)
	fmt.Println("OK")

	var A int
	fmt.Scanf("A is %d", &A)

	r := rand.New(rand.NewSource(42))
	b := r.Intn(p-1) + 1
	B := modexp(g, p, b)
	S := modexp(A, p, b)

	fmt.Println("B is ", B)
	fmt.Println("A is ", A)
	fmt.Println("S is ", S)
}
