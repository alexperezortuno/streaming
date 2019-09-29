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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	base, _ := os.Getwd()
	tmplFile := fmt.Sprintf("%s/src/templates/home.html", base)
	
	tmpl := template.Must(template.ParseFiles(tmplFile))
	_ = tmpl.Execute(w, "")
}

func IndexApiHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	
	_, _ = fmt.Fprintf(w, "Stream app")
}

func StreamHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	mId, err := strconv.Atoi(vars["mId"])

	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	segName, ok := vars["segName"]

	if !ok {
		mediaBase := tools.GetMediaBase(mId)
		m3u8Name := "index.m3u8"
		serveHlsM3u8(response, request, mediaBase, m3u8Name)
	} else {
		mediaBase := tools.GetMediaBase(mId)
		serveHlsTs(response, request, mediaBase, segName)
	}
}

func UploadApiHandler(response http.ResponseWriter, request *http.Request) {
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
	list := request.Form.Get("list")
	
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
			tools.CheckErr(err, 153)
			
			smt, err := db.Prepare("INSERT INTO video (internal_id, name, list) VALUES (?, ?, ?)")
			tools.CheckErr(err, 156)
			
			res, err := smt.Exec(rndString, name, list)
			tools.CheckErr(err, 159)
			
			_, err = res.LastInsertId()
			tools.CheckErr(err, 163)
			defer smt.Close()
			
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

func VideosApiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	base, err := os.Getwd()
	tools.CheckErr(err, 181)

	db, err := sql.Open("sqlite3", base + "/src/stream.db")
	tools.CheckErr(err, 184)
	
	rows, err := db.Query("SELECT name, internal_id AS internalId, list FROM video")
	tools.CheckErr(err, 187)
	
	var name string
	var internalId string
	var list string
	var response structs.VideoResponse
	
	for rows.Next() {
		err = rows.Scan(&name, &internalId, &list)
		response.Message = append(response.Message, structs.Video{Name: name, Id: internalId, List: list})
	}
	
	response.Code = 200
	response.Status = "OK"
	
	_ = rows.Close()
	
	_ = db.Close()
	
	_ = json.NewEncoder(w).Encode(response)
}

func VideoApiHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    vars := mux.Vars(r)
    id := vars["id"]
    var name string
	var internalId string
	var list string
	var response structs.VideoResponse
	var code int16
	var status string
    
    base, err := os.Getwd()
    tools.CheckErr(err, 216)
    
    db, err := sql.Open("sqlite3", base+"/src/stream.db")
    tools.CheckErr(err, 219)
    
	rowsCount, err := db.Query("SELECT count(*) FROM video WHERE internal_id=?", id)
	count := tools.RowCount(rowsCount)
	
	defer rowsCount.Close()
	
	if count <= 0 {
		code = 300
		status = "OK"
	}
    
    if count > 0 {
		rows, err := db.Query("SELECT name, internal_id AS internalId, list FROM video WHERE internal_id=?", id)
		tools.CheckErr(err, 222)
		
		defer rows.Close()
	
		for rows.Next() {
			err = rows.Scan(&name, &internalId, &list)
		
			if err != nil {
				log.Fatal(err.Error())
			} else {
				response.Message = append(response.Message, structs.Video{Name: name, Id: internalId, List: list})
				code = 200
				status = "OK"
			}
		}
	}
    
    response.Code = code
    response.Status = status
	
    _ = db.Close()
	
	defer r.Body.Close()
    
    _ = json.NewEncoder(w).Encode(response)
}