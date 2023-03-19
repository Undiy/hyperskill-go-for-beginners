package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func modexp(g, p, b int) int {
	res := 1
	for i := 0; i < b; i++ {
		res = (res * g) % p
	}
	return res
}

func shift(ch, from, to, n int) rune {
	mod := to - from + 1
	return rune(((ch-from+n)%mod+mod)%mod + from)
}

func shiftChar(ch rune, n int) rune {
	if ch >= 'A' && ch <= 'Z' {
		return shift(int(ch), 'A', 'Z', n)
	} else if ch >= 'a' && ch <= 'z' {
		return shift(int(ch), 'a', 'z', n)
	} else {
		return ch
	}
}

func shiftString(s string, n int) string {
	res := make([]rune, len(s))
	for i, ch := range s {
		res[i] = shiftChar(ch, n)
	}
	return string(res)
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
	fmt.Println(shiftString("Will you marry me?", S))

	reader := bufio.NewReader(os.Stdin)
	ansRaw, _ := reader.ReadString('\n')
	ans := shiftString(ansRaw, -S)
	if ans == "Yeah, okay!\n" {
		fmt.Println(shiftString("Great!", S))
	} else if ans == "Let's be friends.\n" {
		fmt.Println(shiftString("What a pity!", S))
	}
}
