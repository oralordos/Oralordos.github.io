package main

import (
	"crypto/md5"
	"encoding/hex"
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
		http.Error(res, "Server Error", http.StatusInternalServerError)
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
	err = tpl.Execute(res, sortedImages)
	if err != nil {
		http.Error(res, "Server Error", http.StatusInternalServerError)
		return
	}
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
		http.Redirect(res, req, "/admin", http.StatusSeeOther)
		return
	}

	tpl, err := template.ParseFiles("login.gohtml")
	if err != nil {
		http.Error(res, "Server Error", http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(res, failedLogin)
	if err != nil {
		http.Error(res, "Server Error", http.StatusInternalServerError)
		return
	}
}

func logout(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	delete(session.Values, "user")
	store.Save(req, res, session)
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func adminSite(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	if session.Values["user"] != "admin" {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}
	gotFile := false
	if req.Method == "POST" {
		file, _, err := req.FormFile("image")
		if err != nil {
			http.Error(res, "Server Error", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		ckSum := md5.New()
		io.Copy(ckSum, file)
		filename := "images/" + hex.EncodeToString(ckSum.Sum(nil))
		wtr, err := os.Create(filename)
		if err != nil {
			http.Error(res, "Server Error", http.StatusInternalServerError)
			return
		}
		defer wtr.Close()

		_, err = file.Seek(0, 0)
		if err != nil {
			http.Error(res, "Server Error", http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(wtr, file)
		if err != nil {
			http.Error(res, "Server Error", http.StatusInternalServerError)
			return
		}

		gotFile = true
	}

	tpl, err := template.ParseFiles("adminSite.gohtml")
	if err != nil {
		http.Error(res, "Server Error", http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(res, gotFile)
	if err != nil {
		http.Error(res, "Server Error", http.StatusInternalServerError)
		return
	}
}

func getCSS(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "style.css")
}

func toHTTPSHandler(res http.ResponseWriter, req *http.Request) {
	changedURL := "https://" + req.Host[:len(req.Host)-1] + "1/" + req.URL.Path
	http.Redirect(res, req, changedURL, http.StatusSeeOther)
}

func toHTTPHandler(res http.ResponseWriter, req *http.Request) {
	changedURL := "http://" + req.Host[:len(req.Host)-1] + "0/" + req.URL.Path
	http.Redirect(res, req, changedURL, http.StatusSeeOther)
}

func main() {
	imagesHandler := http.StripPrefix("/images/", http.FileServer(http.Dir("images/")))

	http.HandleFunc("/", mainSite)
	http.HandleFunc("/admin", toHTTPSHandler)
	http.Handle("/images/", imagesHandler)
	http.HandleFunc("/login", toHTTPSHandler)
	http.HandleFunc("/logout", toHTTPSHandler)
	http.HandleFunc("/style.css", getCSS)
	go http.ListenAndServe(":9000", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/", toHTTPHandler)
	mux.HandleFunc("/admin/", adminSite)
	mux.Handle("/images/", imagesHandler)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/style.css", getCSS)
	http.ListenAndServeTLS(":9001", "cert.pem", "key.pem", mux)
}
