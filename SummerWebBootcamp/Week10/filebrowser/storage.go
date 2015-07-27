package browser

import (
	"io"
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const delimiter = "/"

type file struct {
	IsFolder bool
	Filename string
	URL      string
}

func getCloudContext(aeCtx context.Context, credentials string) (context.Context, error) {
	conf, err := google.JWTConfigFromJSON(
		[]byte(credentials),
		storage.ScopeFullControl,
	)
	if err != nil {
		return nil, err
	}

	tokenSource := conf.TokenSource(aeCtx)

	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
			Base:   &urlfetch.Transport{Context: aeCtx},
		},
	}
	return cloud.NewContext(appengine.AppID(aeCtx), hc), nil
}

func listFiles(cctx context.Context, bucket, path string) ([]file, error) {
	q := &storage.Query{
		Delimiter: delimiter,
		Prefix:    path,
	}
	objs, err := storage.ListObjects(cctx, bucket, q)
	if err != nil {
		return nil, err
	}
	files := []file{}
	for _, v := range objs.Results {
		f := file{
			URL:      v.Name,
			Filename: v.Name,
			IsFolder: strings.HasSuffix(v.Name, delimiter),
		}
		files = append(files, f)
	}
	return files, nil
}

func getFile(cctx context.Context, bucket, path string) (io.ReadCloser, error) {
	return storage.NewReader(cctx, bucket, path)
}
