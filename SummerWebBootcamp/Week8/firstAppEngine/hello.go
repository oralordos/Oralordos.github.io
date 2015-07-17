package hello

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Hello World!")
}
