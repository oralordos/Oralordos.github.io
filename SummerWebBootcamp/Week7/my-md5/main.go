package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Please pass in a filename to calculate the md5 of")
	}

	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			hash := md5.New()
			io.Copy(hash, f)
			fmt.Printf("%x\t%s\n", hash.Sum(nil), info.Name())
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
}
