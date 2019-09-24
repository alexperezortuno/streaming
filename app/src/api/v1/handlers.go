package v1

import (
	"api/structs"
	"api/tools"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	base, _ := os.Getwd()
	tmplFile := fmt.Sprintf("%s/src/templates/home.html", base)
	
	tmpl := template.Must(template.ParseFiles(tmplFile))
	_ = tmpl.Execute(w, "")
}

func indexApiHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Stream app")
}

func getMediaBase(mId int) string {
	mediaRoot := "media"

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

func serveHlsM3u8(response http.ResponseWriter, request *http.Request, mediaBase, m3u8Name string) {
	base, _ := os.Getwd()
	mediaFile := fmt.Sprintf("%s/src/static/%s/hls/%s", base, mediaBase, m3u8Name)
	http.ServeFile(response, request, mediaFile)
	
	response.Header().Set("Content-Type", "application/x-mpegURL")
}

func serveHlsTs(response http.ResponseWriter, request *http.Request, mediaBase, segName string) {
	base, _ := os.Getwd()
	mediaFile := fmt.Sprintf("%s/src/static/%s/hls/%s", base, mediaBase, segName)
	http.ServeFile(response, request, mediaFile)
	
	response.Header().Set("Content-Type", "video/MP2T")
}

func uploadHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	base, _ := os.Getwd()
	var message []string
	var status string
	var code int16
	fmt.Println("File Upload Endpoint Hit")
	rndString := tools.RandStringBytesMaskImprSrc(rand.NewSource(time.Now().UTC().UnixNano()), 10)
	dirName := "/src/static/media/" + rndString
	
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	_ = request.ParseMultipartForm(10 << 20)
	
	name := request.Form.Get("name")
	
	if name == "" {
		fmt.Println("Error Retrieving the name file")
		return
	}
	
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := request.FormFile("video")
	
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	
	defer file.Close()
	
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	
	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tools.CreateDirIfNotExist(dirName + "/hls")
	
	tempFile, err := ioutil.TempFile(base + dirName + "/hls", "input-*")
	
	if err != nil {
		fmt.Println(err)
	}
	
	defer tempFile.Close()
	
	fileBytes, err := ioutil.ReadAll(file)
	
	if err != nil {
		_ = os.RemoveAll(base + dirName)
		fmt.Println(err)
	}
	
	// read all of the contents of our uploaded file into a
	// byte array
	// write this byte array to our temporary file
	_, err = tempFile.Write(fileBytes)
	
	if err == nil {
		cmd := exec.Command("ffmpeg", "-i", tempFile.Name(), "-profile:v", "baseline", "-level", "3.0", "-s", "1280x720", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", base + dirName + "/hls/" + "index.m3u8")
		randomBytes := &bytes.Buffer{}
		cmd.Stdout = randomBytes
		err = cmd.Start()
		
		if err != nil {
			log.Fatal(err)
			message = []string{"Error to upload file"}
			code = 300
			status = "Error"
		} else {
			db, err := sql.Open("sqlite3", base + "/src/stream.db")
			tools.CheckErr(err)
			
			smt, err := db.Prepare("INSERT INTO video (internal_id, name) SET (?, ?)")
			tools.CheckErr(err)
			
			res, err = smt.Exec(rndString, name)
			tools.CheckErr(err)
			
			id, err := res.LastInsertId()
			checkErr(err)
			
			_ = db.Close()
			
			message = []string{"Successfully Uploaded File"}
			code = 200
			status = "OK"
		}
		
		r := structs.DefaultResponse{Code: code, Status: status, Message: message}
		js, err := json.Marshal(r)
		
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		
		_, _ = response.Write(js)
	}
}
