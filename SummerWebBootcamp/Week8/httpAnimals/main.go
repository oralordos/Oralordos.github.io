package main

import (
	"io"
	"net/http"
)

func catHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, `<!DOCTYPE html>
<html>
  <body>
    <img src="http://lorempixel.com/400/400/cats" alt="Cat Image">
  </body>
</html>`)
}

func dogHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, `<!DOCTYPE html>
<html>
  <body>
    <img src="http://lorempixel.com/400/400/animals/8" alt="Dog Image">
  </body>
</html>`)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cat/", catHandler)
	mux.HandleFunc("/dog/", dogHandler)
	http.ListenAndServe(":9000", mux)
}
