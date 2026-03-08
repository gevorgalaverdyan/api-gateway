package mockservers

import (
	"fmt"
	"net/http"
	"net/url"
)

func RunServers(urls []url.URL) {
	for _, val := range urls {
		val := val
		fmt.Println("Hosting: " + val.Host)
		go func() {
			http.ListenAndServe(":"+val.Port(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello from : %s", r.Host)
			}))
		}()
	}
}
