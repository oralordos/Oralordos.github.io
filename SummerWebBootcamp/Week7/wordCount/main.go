package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func count(r io.Reader) map[string](int) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	words := map[string](int){}
	for scanner.Scan() {
		words[strings.ToLower(scanner.Text())]++
	}
	return words
}

func main() {
	message, _ := os.Open("moby10b.txt")
	words := count(message)
	fmt.Println("Number of the word \"and\":", words["and"])
	fmt.Println("Number of the word \"whale\":", words["whale"])
}
