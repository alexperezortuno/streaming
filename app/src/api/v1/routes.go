package v1

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path"
)

func routes() *mux.Router {
	r := mux.NewRouter()
	base := os.Getenv("BASE")
	// Serve static files
	sf := http.FileServer(http.Dir(path.Join(base, "app/src/static")))
	mf := http.FileServer(http.Dir(path.Join(base, "app/src/static/media")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", sf))
	r.PathPrefix("/media/").Handler(http.StripPrefix("/media/", mf))

	r.HandleFunc("/api/v1", indexHandler).Methods("GET")
	r.HandleFunc("/api/media/{mId:[0-9]+}/stream/", streamHandler).Methods("GET")
	r.HandleFunc("/api/media/{mId:[0-9]+}/stream/{segName:index[0-9]+.ts}", streamHandler).Methods("GET")

	return r
}
