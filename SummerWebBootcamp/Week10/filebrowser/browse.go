package browser

import (
	"io"
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/julienschmidt/httprouter"
)

type browseModel struct {
	Path   string
	Bucket string
	Files  []file
}

func handlePath(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(req)
	s := getSession(ctx, req)
	if s.Bucket == "" {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	cctx, err := getCloudContext(ctx, s.Credentials)
	if err != nil {
		http.Error(res, "Server Error", http.StatusInternalServerError)
		log.Errorf(ctx, err.Error())
		return
	}

	path := p.ByName("path")[1:]

	if !strings.HasSuffix(path, delimiter) && path != "" {
		io.WriteString(res, path)
	} else {
		files, err := listFiles(cctx, s.Bucket, path)
		if err != nil {
			http.Error(res, "Server Error", http.StatusInternalServerError)
			log.Errorf(ctx, err.Error())
			return
		}
		data := browseModel{
			Path:   "/" + path,
			Bucket: s.Bucket,
			Files:  files,
		}
		err = tpl.ExecuteTemplate(res, "browse", data)
		if err != nil {
			http.Error(res, "Server Error", http.StatusInternalServerError)
			log.Errorf(ctx, err.Error())
			return
		}
	}
}
