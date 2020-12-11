package tools

import (
	"../structs"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	chars         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func WaitForShutdown(server *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_ = server.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}

func RandStringBytesMaskImprSrc(src rand.Source, n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(chars) {
			b[i] = chars[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func CreateDirIfNotExist(dir string) {
	base, _ := os.Getwd()
	dir = base + "/" + dir

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)

		if err != nil {
			panic(err)
		}
	}
}

func CheckErr(err error, line int) {
	if err != nil {
		log.Fatal("Error DB: " + err.Error() + " line: " + strconv.Itoa(line))
		panic(err)
	}
}

func GetMediaBase(mId int) string {
	mediaRoot := "media"

	return fmt.Sprintf("%s/%d", mediaRoot, mId)
}

func RowCount(rows *sql.Rows) (count int) {

	for rows.Next() {
		err := rows.Scan(&count)
		CheckErr(err, 88)
	}

	return count
}

func DBConnection() *sql.DB {
	base, err := os.Getwd()
	CheckErr(err, 216)

	db, err := sql.Open("sqlite3", base+"/stream.db")
	CheckErr(err, 219)

	return db
}

func JsonResponse(w http.ResponseWriter, status int, results structs.DefaultResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(results)
}
