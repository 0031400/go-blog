package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	var err error
	loalConfig()
	initDB()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !authMiddler(w, r) {
			return
		}
		corsMiddler(w, r)
		path := r.URL.Path
		method := r.Method
		switch {
		case method == http.MethodOptions:
			w.WriteHeader(http.StatusNoContent)
		case path == "/admin/post" && method == http.MethodPut:
			err = handlerPostNew(w, r)
		case path == "/admin/post" && method == http.MethodDelete:
			err = handlerPostDelete(w, r)
		case path == "/admin/post" && method == http.MethodPost:
			err = handlerPostUpdate(w, r)
		case path == "/admin/tag" && method == http.MethodPut:
			err = handlerTagNew(w, r)
		case path == "/admin/tag" && method == http.MethodDelete:
			err = handlerTagDelete(w, r)
		case path == "/admin/tag" && method == http.MethodPost:
			err = handlerTagUpdate(w, r)
		case path == "/admin/category" && method == http.MethodPut:
			err = handlerCategoryNew(w, r)
		case path == "/admin/category" && method == http.MethodDelete:
			err = handlerCategoryDelete(w, r)
		case path == "/admin/category" && method == http.MethodPost:
			err = handlerCategoryUpdate(w, r)
		case path == "/post/list" && method == http.MethodGet:
			err = handlerPostList(w, r)
		case path == "/post" && method == http.MethodGet:
			err = handlerPostDetail(w, r)
		case path == "/tag" && method == http.MethodGet:
			err = handlerTagList(w, r)
		case path == "/category" && method == http.MethodGet:
			err = handlerCategoryList(w, r)
		default:
			http.Error(w, "not found", 404)
		}
		if err != nil {
			log.Println(err)
		}
	})
	fmt.Printf("The server is running on %s\n", config.ListenAddr)
	err = http.ListenAndServe(config.ListenAddr, nil)
	if err != nil {
		panic(err)
	}
}
func authMiddler(w http.ResponseWriter, r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/admin") && r.Header.Get("Authorization") != base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", config.AdminAccount, config.AdminPassword))) {
		http.Error(w, "auth fail", 403)
		return false
	}
	return true
}
func corsMiddler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
