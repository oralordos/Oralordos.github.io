package twitter

import (
	"html/template"
	"net/http"
	"time"

	"golang.org/x/net/context"

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
	Username   string `datastore:"-"`
	Message    string
	SubmitTime time.Time
}

type mainpageData struct {
	Tweets   []tweet
	Logged   bool
	Username string
}

type profileData struct {
	Tweets  []tweet
	Profile profile
}

type loginData struct {
	ErrorMessage string
	Username     string
}

const (
	minUsernameSize = 5
	maxUsernameSize = 20
	loginDuration   = 60 * 60 * 24 // 1 Day
)

var tpl = template.New("templates")

func init() {
	_, err := tpl.ParseGlob("templates/*.gohtml")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", handle)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/CreateProfile", createProfile)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
}

func confirmCreateProfile(ctx context.Context, username string) bool {
	_, err := getProfileByUsername(ctx, username)
	return len(username) >= minUsernameSize && len(username) <= maxUsernameSize &&
		err == datastore.ErrNoSuchEntity
}

func handleLogin(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	u := user.Current(ctx)

	cookie, err := req.Cookie("login")
	if err != http.ErrNoCookie {
		http.Redirect(res, req, "/"+cookie.Value, http.StatusSeeOther)
		return
	}

	currentProfile, err := getProfileByEmail(ctx, u.Email)
	if err == datastore.ErrNoSuchEntity {
		http.Redirect(res, req, "/CreateProfile", http.StatusSeeOther)
		return
	} else if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Get profile error: %s\n", err.Error())
		return
	}

	login := loginData{
		Username: currentProfile.Username,
	}
	if req.Method == "POST" {
		username := req.FormValue("username")
		p, err := getProfileByUsername(ctx, username)
		if err != nil {
			login.ErrorMessage = "No such username"
		} else if p.Email != u.Email {
			login.ErrorMessage = "Not your profile"
		} else {
			c := http.Cookie{
				Name:   "login",
				Value:  username,
				MaxAge: loginDuration,
			}
			http.SetCookie(res, &c)
			http.Redirect(res, req, "/"+username, http.StatusSeeOther)
			return
		}
	}
	err = tpl.ExecuteTemplate(res, "login.gohtml", login)
	if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Template Execute Error: %s\n", err.Error())
		return
	}
}

func handleLogout(res http.ResponseWriter, req *http.Request) {
	http.SetCookie(res, &http.Cookie{Name: "login", MaxAge: -1})
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func createProfile(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	u := user.Current(ctx)

	if req.Method == "POST" {
		username := req.FormValue("username")
		if !confirmCreateProfile(ctx, username) {
			http.Error(res, "Invalid input!", http.StatusBadRequest)
			log.Warningf(ctx, "Invalid profile information from %s\n", req.RemoteAddr)
			return
		}
		key := datastore.NewKey(ctx, "profile", u.Email, 0, nil)
		p := profile{
			Username: username,
			Email:    u.Email,
		}
		http.SetCookie(res, &http.Cookie{Name: "login", Value: username, MaxAge: loginDuration})
		_, err := datastore.Put(ctx, key, &p)
		if err != nil {
			http.Error(res, "Server error!", http.StatusInternalServerError)
			log.Errorf(ctx, "Create profile Error: %s\n", err.Error())
			return
		}
	}

	_, err := getProfileByEmail(ctx, u.Email)
	if err == nil {
		http.Redirect(res, req, "login", http.StatusSeeOther)
		return
	} else if err != datastore.ErrNoSuchEntity {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Get profile Error: %s\n", err.Error())
		return
	}

	err = tpl.ExecuteTemplate(res, "createProfile.gohtml", nil)
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
	data := mainpageData{
		Tweets: tweets,
	}

	c, err := req.Cookie("login")
	if err == nil {
		data.Logged = true
		data.Username = c.Value
	} else {
		data.Logged = false
	}

	err = tpl.ExecuteTemplate(res, "index.gohtml", data)
	if err != nil {
		http.Error(res, "Server error!", http.StatusInternalServerError)
		log.Errorf(ctx, "Template Execute Error: %s\n", err.Error())
		return
	}
}
