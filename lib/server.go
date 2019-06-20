package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const retryAfterInputFormat =  "2006-01-02 15:04:05-0700"
const retryAfterFormat = time.RFC1123

func readFile(filename string) []byte {
	if data, err := ioutil.ReadFile(filename); err != nil {
		panic(err)
	} else {
		return data
	}
}

func getRetryAfter() *time.Time {
	var retryAfter = os.Getenv("RETRY_AFTER")
	if retryAfter == "" {
		return nil
	}
	if t, err := time.Parse(retryAfterInputFormat, retryAfter); err != nil {
		panic(err)
	} else {
		return &t
	}
}

func getJson() []byte {
	return readFile("503.json")
}

func getHtml() []byte {
	return readFile("503.html")
}

var (
	json = getJson()
	html = getHtml()
)

// Accept ヘッダが json もしくは path が json で終わる場合
func isExpectJsonRequest(r *http.Request) bool {
	if strings.Contains(r.Header.Get("Accept"), "/json") {
		return true
	}
	if strings.HasSuffix(r.URL.Path, ".json") {
		return true
	}
	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request URL: %s%s\n", r.Host, r.URL.String())

	var contentType = "text/html; charset=utf-8"
	var body = html

	// json レスポンス
	if isExpectJsonRequest(r) {
		contentType = "application/json; charset=utf-8"
		body = json
	}

	w.Header().Add("Content-Type", contentType)

	// Retry-After が指定されていたらヘッダ追加
	if retryAfter := getRetryAfter(); retryAfter != nil {
		w.Header().Add("Retry-After", retryAfter.Format(retryAfterFormat))
	}

	// ステータスコードは 503 固定
	w.WriteHeader(http.StatusServiceUnavailable)

	w.Write(body)
}

func CreateServer() *http.Server {
	mux := http.NewServeMux()
	// assets
	mux.Handle("/assets/", http.FileServer(http.Dir("./")))
	mux.HandleFunc("/", handler)
	return &http.Server{
		Handler: mux,
		ReadTimeout: 10  * time.Second,
		WriteTimeout: 10  * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
