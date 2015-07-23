package chat

import "net/http"

func init() {
	http.Handle("/", http.FileServer(http.Dir("public/")))
	http.Handle("/api/", newAPI("/api/"))
	http.HandleFunc("/_ah/channel/connected", connectClient)
	http.HandleFunc("/_ah/channel/disconnected", disconnectClient)
}
