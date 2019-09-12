package v1

import (
    "github.com/gorilla/mux"
)

func routes() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/api/v1", indexHandler).Methods("GET")
    r.HandleFunc("/api/v1/media/{mId:[0-9]+}/stream/", streamHandler).Methods("GET")
    r.HandleFunc("/api/v1/media/{mId:[0-9]+}/stream/{segName:index[0-9]+.ts}", streamHandler).Methods("GET")
    
    return r
}
