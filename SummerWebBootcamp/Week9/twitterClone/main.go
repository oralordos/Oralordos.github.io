package twitter

import (
	"html/template"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

type tweet struct {
	Message    []string
	SubmitTime time.Time
}

type mainpageData struct {
	Tweets []tweet
	Logged bool
	Email  string
}

func init() {
	http.HandleFunc("/", handle)
}

func handle(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	u := user.Current(ctx)

	// Get recent tweets
	query := datastore.NewQuery("Tweets")
	tweets := []tweet{}
	_, err := query.GetAll(ctx, &tweets)
	if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Query Error: %s\n", err.Error())
		return
	}

	// Create template
	data := mainpageData{
		Tweets: tweets,
		Logged: u != nil,
	}
	if data.Logged {
		data.Email = u.Email
	}

	tpl, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Template Parse Error: %s\n", err.Error())
		return
	}
	err = tpl.Execute(res, data)
	if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Template Execute Error: %s\n", err.Error())
		return
	}
}
