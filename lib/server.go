package lib

import (
	"log"
	"net/http"
	"strings"
	"time"
)

const retryAfterFormat = time.RFC1123

// Accept ヘッダが Json もしくは path が Json で終わる場合
func isExpectJsonRequest(r *http.Request) bool {
	if strings.Contains(strings.ToLower(r.Header.Get("Accept")), "/json") {
		return true
	}
	if strings.HasSuffix(strings.ToLower(r.URL.Path), ".json") {
		return true
	}
	return false
}

func createHandler(contents *ResponseContents, scheme string) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URL: %s://%s%s\n", scheme, r.Host, r.URL.String())

		var contentType = "text/html; charset=utf-8"
		var body = contents.html

		// Json レスポンス
		if isExpectJsonRequest(r) {
			contentType = "application/json; charset=utf-8"
			body = contents.json
		}

		w.Header().Add("Content-Type", contentType)

		// Retry-After が指定されていたらヘッダ追加
		if contents.retryAfter != nil {
			w.Header().Add("Retry-After", contents.retryAfter.Format(retryAfterFormat))
		}

		// ステータスコードは 503 固定
		w.WriteHeader(http.StatusServiceUnavailable)

		w.Write(body)
	}
}

func CreateServer(contents *ResponseContents, scheme string) *http.Server {
	mux := http.NewServeMux()
	// assets
	mux.Handle("/assets/", http.FileServer(http.Dir("./")))
	mux.HandleFunc("/", createHandler(contents, scheme))
	return &http.Server{
		Handler: mux,
		ReadTimeout: 10  * time.Second,
		WriteTimeout: 10  * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
