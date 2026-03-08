package main

import (
	"fmt"
	mockservers "gateway/mock_servers"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	routes := map[string]url.URL{
		"/users":    {Scheme: "http", Host: "localhost:9000"},
		"/products": {Scheme: "http", Host: "localhost:9001"},
		"/services": {Scheme: "http", Host: "localhost:9002"},
	}

	hosts := []url.URL{}
	for _, h := range routes {
		hosts = append(hosts, h)
	}

	mockservers.RunServers(hosts)

	for rPath, target := range routes {
		http.HandleFunc(rPath, func(w http.ResponseWriter, r *http.Request) {
			// fmt.Fprintf(w, "Forwarding to: %s", target)
			proxy := httputil.NewSingleHostReverseProxy(&target)
			proxy.ServeHTTP(w, r)
		})
	}

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
