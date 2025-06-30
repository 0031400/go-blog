package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	var err error
	initDB()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case path == "/admin/post" && r.Method == http.MethodPut:
			handlerPostNew(w, r)
		case strings.HasPrefix(path, "/admin/post") && r.Method == http.MethodDelete:
			handlerPostDelete(w, r)
		case path == "/admin/post" && r.Method == http.MethodPost:
			handlerPostUpdate(w, r)
		default:
			http.Error(w, "not found", 404)
		}
	})
	fmt.Print("The server is running on 127.0.0.1:8080\n")
	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		panic(err)
	}
}
