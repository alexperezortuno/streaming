package v1

import (
    "fmt"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "strconv"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
    log.Print("index")
    http.ServeFile(w, r, "index.html")
}

func getMediaBase(mId int) string {
    mediaRoot := "assets/media"
    
    return fmt.Sprintf("%s/%d", mediaRoot, mId)
}

func streamHandler(response http.ResponseWriter, request *http.Request) {
    vars := mux.Vars(request)
    mId, err := strconv.Atoi(vars["mId"])
    
    if err != nil {
        response.WriteHeader(http.StatusNotFound)
        
        return
    }
    
    segName, ok := vars["segName"]
    
    if !ok {
        mediaBase := getMediaBase(mId)
        m3u8Name := "index.m3u8"
        serveHlsM3u8(response, request, mediaBase, m3u8Name)
    } else {
        mediaBase := getMediaBase(mId)
        serveHlsTs(response, request, mediaBase, segName)
    }
}

func serveHlsM3u8(w http.ResponseWriter, r *http.Request, mediaBase, m3u8Name string) {
    mediaFile := fmt.Sprintf("%s/hls/%s", mediaBase, m3u8Name)
    http.ServeFile(w, r, mediaFile)
    
    w.Header().Set("Content-Type", "application/x-mpegURL")
}

func serveHlsTs(w http.ResponseWriter, r *http.Request, mediaBase, segName string) {
    mediaFile := fmt.Sprintf("%s/hls/%s", mediaBase, segName)
    http.ServeFile(w, r, mediaFile)
    
    w.Header().Set("Content-Type", "video/MP2T")
}