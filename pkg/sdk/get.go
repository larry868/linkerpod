package sdk

import (
	"io"
	"log"
	"net/http"
	"time"
)

// ApiGet returns the decoded json and a success flag
func ApiGet(host string, apiname string, target string, ferror func(statuscode int, target string, err error)) (body []byte, httpStatus int, start time.Time) {

	start = time.Now()

	// call the server
	var status int
	resp, errGet := http.Get(host + "/api/" + apiname)
	if resp != nil {
		status = resp.StatusCode
	}
	if errGet != nil {
		log.Println(errGet)
		if ferror != nil {
			ferror(status, target, errGet)
			return
		}
		return []byte{}, status, start
	}
	if status != http.StatusOK {
		return []byte{}, status, start
	}

	// extract the response
	var errRead error
	body, errRead = io.ReadAll(resp.Body)
	if errRead != nil {
		log.Println(errRead)
		if ferror != nil {
			ferror(status, target, errRead)
			return
		}
		return []byte{}, status, start
	}

	return body, status, start
}
