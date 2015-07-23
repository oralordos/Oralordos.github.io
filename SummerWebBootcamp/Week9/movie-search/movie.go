package search

import (
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type movie struct {
	Name     string
	URL      string
	Summary  string
	ImageURL string
}

func handleAdd(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	var err error

	if req.Method == "POST" {
		name := strings.TrimSpace(req.FormValue("name"))
		summary := req.FormValue("summary")
		imageURL := req.FormValue("imageURL")

		mov := &movie{
			Name:     name,
			Summary:  summary,
			URL:      strings.ToLower(strings.Replace(name, " ", "", -1)),
			ImageURL: imageURL,
		}
		err = addMovie(ctx, mov)
		if err != nil {
			http.Error(res, "Server error", http.StatusInternalServerError)
			log.Errorf(ctx, "%v\n", err)
			return
		}
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	err = tpl.ExecuteTemplate(res, "addMovie", nil)
	if err != nil {
		http.Error(res, "Server error", http.StatusInternalServerError)
		log.Errorf(ctx, "%v\n", err)
		return
	}
}
