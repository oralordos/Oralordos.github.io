package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("asdf"))

type fileTimes []time.Time

func (ft *fileTimes) Len() int {
	return len(*ft)
}

func (ft *fileTimes) Less(i, j int) bool {
	return (*ft)[i].After((*ft)[j])
}

func (ft *fileTimes) Swap(i, j int) {
	(*ft)[i], (*ft)[j] = (*ft)[j], (*ft)[i]
}

func mainSite(res http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles("mainSite.gohtml")
	if err != nil {
		http.Error(res, "Server Error", 500)
		return
	}
	times := fileTimes{}
	images := map[time.Time]string{}
	filepath.Walk("images", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileTime := info.ModTime()
			images[fileTime] = filepath.ToSlash("/" + path)
			times = append(times, fileTime)
		}
		return nil
	})
	sort.Sort(&times)
	sortedImages := make([]string, len(times))
	for i, v := range times {
		sortedImages[i] = images[v]
	}
	tpl.Execute(res, sortedImages)
}

func login(res http.ResponseWriter, req *http.Request) {
	failedLogin := false
	session, _ := store.Get(req, "session")
	if req.Method == "POST" {
		username := req.FormValue("username")
		password := req.FormValue("password")
		if username == "me" && password == "you" {
			session.Values["user"] = "admin"
			store.Save(req, res, session)
		} else {
			failedLogin = true
		}
	}
	if session.Values["user"] == "admin" {
		http.Redirect(res, req, "/admin", 303)
		return
	}

	tpl, err := template.ParseFiles("login.gohtml")
	if err != nil {
		http.Error(res, "Server Error", 500)
		return
	}
	tpl.Execute(res, failedLogin)
}

func logout(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	delete(session.Values, "user")
	store.Save(req, res, session)
	http.Redirect(res, req, "/", 303)
}

func adminSite(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	if session.Values["user"] != "admin" {
		http.Redirect(res, req, "/login", 303)
		return
	}
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

func toHTTPSHandler(res http.ResponseWriter, req *http.Request) {
	changedURL := "https://" + req.Host[:len(req.Host)-1] + "1/" + req.URL.Path
	http.Redirect(res, req, changedURL, 303)
}

func toHTTPHandler(res http.ResponseWriter, req *http.Request) {
	changedURL := "http://" + req.Host[:len(req.Host)-1] + "0/" + req.URL.Path
	http.Redirect(res, req, changedURL, 303)
}

func main() {
	imagesHandler := http.StripPrefix("/images/", http.FileServer(http.Dir("images/")))

	http.HandleFunc("/", mainSite)
	http.HandleFunc("/admin", toHTTPSHandler)
	http.Handle("/images/", imagesHandler)
	http.HandleFunc("/login", toHTTPSHandler)
	http.HandleFunc("/logout", toHTTPSHandler)
	go http.ListenAndServe(":9000", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/", toHTTPHandler)
	mux.HandleFunc("/admin/", adminSite)
	mux.Handle("/images/", imagesHandler)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	http.ListenAndServeTLS(":9001", "cert.pem", "key.pem", mux)
}
