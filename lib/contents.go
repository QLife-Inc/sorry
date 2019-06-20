package lib

import (
	"io/ioutil"
	"os"
	"time"
)

const retryAfterInputFormat =  "2006-01-02 15:04:05-0700"

type ResponseContents struct {
	html       []byte
	json       []byte
	retryAfter *time.Time
}

func NewResponseContents() (*ResponseContents, error) {
	html, err := getHtml()
	if err != nil {
		return nil, err
	}
	json, err := getJson()
	if err != nil {
		return nil, err
	}
	retryAfter, err := getRetryAfter()
	if err != nil {
		return nil, err
	}
	return &ResponseContents{html: html, json: json, retryAfter: retryAfter}, nil
}

func readFile(filename string) ([]byte, error) {
	if data, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func getRetryAfter() (*time.Time, error) {
	var retryAfter = os.Getenv("RETRY_AFTER")
	if retryAfter == "" {
		return nil, nil
	}
	if t, err := time.Parse(retryAfterInputFormat, retryAfter); err != nil {
		return nil, err
	} else {
		return &t, nil
	}
}

func getJson() ([]byte, error) {
	return readFile("503.json")
}

func getHtml() ([]byte, error) {
	return readFile("503.html")
}
