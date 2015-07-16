package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
)

func mainSite(res http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles("mainSite.gohtml")
	if err != nil {
		http.Error(res, "Server Error", 500)
		return
	}
	tpl.Execute(res, []string{
		"https://Oralordos.github.io/images/background-213649_1280.jpg",
		"https://Oralordos.github.io/images/code-113611_1280.jpg",
	})
}

func login(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {

	}
	http.ServeFile(res, req, "login.html")
}

func adminSite(res http.ResponseWriter, req *http.Request) {
	gotFile := false
	if req.Method == "POST" {
		file, header, err := req.FormFile("image")
		if err != nil {
			http.Error(res, "Server Error", 500)
			return
		}
		defer file.Close()

		filename := "images/" + header.Filename
		wtr, err := os.Create(filename)
		if err != nil {
			http.Error(res, "Server Error", 500)
			return
		}
		defer wtr.Close()

		_, err = io.Copy(wtr, file)
		if err != nil {
			http.Error(res, "Server Error", 500)
			return
		}

		gotFile = true
	}

	tpl, err := template.ParseFiles("adminSite.gohtml")
	if err != nil {
		http.Error(res, "Server Error", 500)
		return
	}
	tpl.Execute(res, gotFile)
}

func main() {
	http.HandleFunc("/", mainSite)
	go http.ListenAndServe(":9000", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/admin", adminSite)
	// mux.HandleFunc("/login", login)
	http.ListenAndServeTLS(":9001", "cert.pem", "key.pem", mux)
}
