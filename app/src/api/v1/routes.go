package v1

import (
	"api/structs"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path"
)

var routes = structs.Routes{
	structs.Route{
		"Index",
		"GET",
		"/",
		IndexHandler,
	},
	structs.Route{
		"IndexHome",
		"GET",
		"/api",
		IndexApiHandler,
	},
	structs.Route{
		"Upload",
		"POST",
		"/api/upload",
		UploadApiHandler,
	},
	structs.Route{
		"Video",
		"GET",
		"/api/video",
		VideosApiHandler,
	},
	structs.Route{
		"VideoDetail",
		"GET",
		"/api/video/{id:[A-z0-9]+}",
		VideoApiHandler,
	},
	structs.Route{
		"Media",
		"GET",
		"/api/media/{mId:[A-z0-9]+}/stream/",
		StreamHandler,
	},
	structs.Route{
		"MediaDetil",
		"GET",
		"/api/media/{mId:[A-z0-9]+}/stream/",
		StreamHandler,
	},
}

func NewRoutes() *mux.Router {
		r := mux.NewRouter().StrictSlash(true)
	base, _ := os.Getwd()
	// Serve static files
	sf := http.FileServer(http.Dir(path.Join(base, "src/static")))
	mf := http.FileServer(http.Dir(path.Join(base, "src/static/media")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", sf))
	r.PathPrefix("/media/").Handler(http.StripPrefix("/media/", mf))
	
	for _, route := range routes {
		r.
			Name(route.Name).
			Methods(route.Method).
			Path(route.Pattern).
			Handler(route.HandleFunc)
	}

	return r
}
