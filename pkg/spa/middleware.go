package spa

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// logger is a middleware handler that does request logging
// https://drstearns.github.io/tutorials/gomiddleware/
type logger struct {
	http.Handler
	apicalls int
}

// newLogger constructs a new Logger middleware handler
func newLogger(handlerToWrap http.Handler) *logger {
	return &logger{handlerToWrap, 0}
}

// ServeHTTP handles the request by passing it to the real
// handler and logging the request details
func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if strings.HasPrefix(r.URL.Path, "/api") {
		l.apicalls++
		call := l.apicalls
		strr := r.Method + " " + r.RequestURI
		fmt.Printf(">>API call \033[7m[#%v[\033[0m %s >>\n", call, strr)
		l.Handler.ServeHTTP(w, r)
		fmt.Printf("<<API call \033[7m]#%v]\033[0m %s << %s\n", call, strr, time.Since(start))
	} else {
		l.Handler.ServeHTTP(w, r)
		fmt.Printf(">>HTTP request: %s %s ✓\n", r.Method, r.RequestURI)
	}
}

// middlewareNoCache add no-cach in every responses' header
func middlewareNoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}
