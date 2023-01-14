package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// ServeApiHealth
func ServeHealth() func(http.ResponseWriter, *http.Request) {
	counter := 0
	return func(w http.ResponseWriter, r *http.Request) {
		counter++
		json.NewEncoder(w).Encode(map[string]string{"health": "live", "counter": strconv.Itoa(counter)})
	}
}
