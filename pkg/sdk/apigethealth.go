package sdk

import (
	"encoding/json"
	"log"
	"net/http"
)

// ApiGetHealth issues a GET requests to /api/health
// returning a text with a server call counter.
// This is necessarilly done async in a seperate go routine, see https://golang.org/pkg/syscall/js/#FuncOf
func ApiGetHealth() (success bool) {
	result := make(chan bool, 1)
	go func(result chan bool) {

		body, status := ApiGet("", "health")
		if status != http.StatusOK {
			log.Println("health: dead")
			result <- false
			return
		}

		// parse received data
		type dataHealth struct {
			Health  string `json:"health"`
			Counter string `json:"counter"`
		}
		data := dataHealth{}
		jsonErr := json.Unmarshal(body, &data)
		if jsonErr != nil {
			log.Printf("health: %s\n", jsonErr.Error())
			result <- false
			return
		}

		// display result
		log.Printf("%v, call counter: %s", data.Health, data.Counter)
		result <- true
	}(result)

	success = <-result
	close(result)
	return success
}
