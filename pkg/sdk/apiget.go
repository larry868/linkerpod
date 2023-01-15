package sdk

import (
	"io"
	"log"
	"net/http"
)

// ApiGet issues a GET request and returns the decoded json and a success flag
func ApiGet(host string, apiname string) (body []byte, httpStatus int) {

	// call the server
	var status int
	resp, errGet := http.Get(host + "/api/" + apiname)
	if resp != nil {
		status = resp.StatusCode
	}
	if status != http.StatusOK || errGet != nil {
		if errGet != nil {
			log.Println(errGet)
		}
		return []byte{}, status
	}

	// extract the response
	var errRead error
	body, errRead = io.ReadAll(resp.Body)
	if errRead != nil {
		log.Println(errRead)
		return []byte{}, status
	}

	return body, status
}
