package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Not enough command line arguments")
	}

	readFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer readFile.Close()

	bytes, err := ioutil.ReadAll(readFile)
	if err != nil {
		log.Fatalln(err)
	}

	if len(os.Args) >= 2 {
		writeFile, err := os.Create(os.Args[2])
		if err != nil {
			log.Fatalln(err)
		}
		defer writeFile.Close()
		_, err = writeFile.Write(bytes)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		str := string(bytes)
		fmt.Println(str)
	}
}
