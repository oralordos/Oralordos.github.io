package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Please pass in a filename to calculate the md5 of")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer f.Close()

	hash := md5.New()
	io.Copy(hash, f)
	fmt.Printf("%x\n", hash.Sum(nil))
}
