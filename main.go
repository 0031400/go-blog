package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	var err error
	initDB()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method
		switch {
		case path == "/admin/post" && method == http.MethodPut:
			handlerPostNew(w, r)
		case path == "/admin/post/" && method == http.MethodDelete:
			handlerPostDelete(w, r)
		case path == "/admin/post" && method == http.MethodPost:
			handlerPostUpdate(w, r)
		case path == "/admin/tag" && method == http.MethodPut:
			err = handlerTagNew(w, r)
		case path == "/admin/tag/" && method == http.MethodDelete:
			err = handlerTagDelete(w, r)
		case path == "/admin/tag" && method == http.MethodPost:
			err = handlerTagUpdate(w, r)
		case path == "/admin/category" && method == http.MethodPut:
			err = handlerCategoryNew(w, r)
		case path == "/admin/category/" && method == http.MethodDelete:
			err = handlerCategoryDelete(w, r)
		case path == "/admin/category" && method == http.MethodPost:
			err = handlerCategoryUpdate(w, r)
		case path == "/post/list" && method == http.MethodGet:
			err = handlerPostList(w, r)
		// case path == "/post" && method == http.MethodGet:
		// 	err = handlerPostDetail(w, r)
		// case path == "/tag" && method == http.MethodGet:
		// 	err = handlerTagList(w, r)
		// case path == "/category" && method == http.MethodGet:
		// 	err = handlerCategoryList(w, r)
		default:
			http.Error(w, "not found", 404)
		}
		if err != nil {
			log.Println(err)
		}
	})
	fmt.Print("The server is running on 127.0.0.1:8080\n")
	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		panic(err)
	}
}
