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

type profile struct {
	Username string
	Email    string
}

type tweet struct {
	Message    string
	SubmitTime time.Time
}

type mainpageData struct {
	Tweets   []tweet
	Logged   bool
	Email    string
	LoginURL string
}

type profileData struct {
	Tweets  []tweet
	Profile profile
}

var tpl = template.New("templates")

func init() {
	_, err := tpl.ParseFiles("templates/index.gohtml", "templates/createProfile.gohtml", "templates/profile.gohtml")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", handle)
	http.HandleFunc("/CreateProfile", createProfile)
}

func confirmCreateProfile(username string) bool {
	return len(username) > 5
}

func createProfile(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	u := user.Current(ctx)

	if req.Method == "POST" {
		username := req.FormValue("username")
		if !confirmCreateProfile(username) {
			http.Error(res, "Invalid input!", http.StatusBadRequest)
			log.Warningf(ctx, "Invalid profile information from %s\n", req.RemoteAddr)
			return
		}
		// TODO Make sure username is not taken
		key := datastore.NewKey(ctx, "profile", u.Email, 0, nil)
		p := profile{
			Username: username,
			Email:    u.Email,
		}
		_, err := datastore.Put(ctx, key, &p)
		if err != nil {
			http.Error(res, "Server error!", http.StatusInternalServerError)
			log.Errorf(ctx, "Create profile Error: %s\n", err.Error())
			return
		}
	}
	err := tpl.ExecuteTemplate(res, "createProfile.gohtml", nil)
	if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Template Execute Error: %s\n", err.Error())
		return
	}
}

func getProfile(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	username := req.URL.Path[1:]

	p, err := getProfileByUsername(ctx, username)
	if err == datastore.ErrNoSuchEntity {
		http.NotFound(res, req)
		return
	} else if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Get Profile Error: %s\n", username)
		return
	}

	if req.Method == "POST" {
		u := user.Current(ctx)
		message := req.FormValue("message")
		if p.Email != u.Email {
			http.Error(res, "Unauthorized post", http.StatusUnauthorized)
			return
		}
		t := tweet{
			Message:    message,
			SubmitTime: time.Now(),
		}
		err := postTweet(ctx, &t, u.Email)
		if err != nil {
			http.Error(res, "Server error!", http.StatusInternalServerError)
			log.Errorf(ctx, "Put Tweet Error: %s\n", err.Error())
			return
		}
	}

	tweets, err := getTweets(ctx, p.Email)
	if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Query Error: %s\n", err.Error())
		return
	}

	pd := profileData{
		Tweets:  tweets,
		Profile: *p,
	}

	err = tpl.ExecuteTemplate(res, "profile.gohtml", pd)
	if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Template Execute Error: %s\n", err.Error())
		return
	}
}

func handle(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		getProfile(res, req)
		return
	}

	ctx := appengine.NewContext(req)
	u := user.Current(ctx)

	if u != nil {
		_, err := getProfileByEmail(ctx, u.Email)
		if err == datastore.ErrNoSuchEntity {
			http.Redirect(res, req, "/CreateProfile", http.StatusSeeOther)
			return
		} else if err != nil {
			http.Error(res, "Server error!", http.StatusInternalServerError)
			log.Errorf(ctx, "Datastore get Error: %s\n", err.Error())
			return
		}
	}

	// Get recent tweets
	tweets, err := getTweets(ctx, "")
	if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Query Error: %s\n", err.Error())
		return
	}

	// Create template
	loginURL, err := user.LoginURL(ctx, "/")
	if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Login URL Error: %s\n", err.Error())
		return
	}
	data := mainpageData{
		Tweets:   tweets,
		Logged:   u != nil,
		LoginURL: loginURL,
	}
	if data.Logged {
		data.Email = u.Email
	}

	err = tpl.ExecuteTemplate(res, "index.gohtml", data)
	if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Template Execute Error: %s\n", err.Error())
		return
	}
}
