package todo

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/", handle)
	http.HandleFunc("/todo.json", jsonServe)
}

func handle(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}

type list struct {
	Test string
}

func jsonServe(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		saveJSON(res, req)
	} else {
		getJSON(res, req)
	}
}

func saveJSON(res http.ResponseWriter, req *http.Request) {
	// TODO
}

func getJSON(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	u := user.Current(ctx)
	key := datastore.NewKey(ctx, "List", u.Email, 0, nil)
	var l list
	err := datastore.Get(ctx, key, &l)
	if err != nil {
		http.Error(res, "Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}
	enc := json.NewEncoder(res)
	enc.Encode(l)
}
