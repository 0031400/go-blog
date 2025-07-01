package main

import (
	"fmt"
	"log"
	"net/http"
)

func handlerPostNew(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	date := r.Form.Get("date")
	content := r.Form.Get("content")
	brief := r.Form.Get("brief")
	requestTagUUIDs := r.Form.Get("tags")
	categoryUUID := r.Form.Get("category")
	tagUUIDs, err := parseStrings(requestTagUUIDs)
	if err != nil {
		log.Println(err)
		http.Error(w, "fail to parse tags", 500)
		return
	}
	if !verifyDate(date) {
		http.Error(w, "date not verified", 400)
		return
	}
	err = postNew(title, date, brief, content, tagUUIDs, categoryUUID)
	if err != nil {
		log.Println(err)
		http.Error(w, "fail to new post", 500)
		return
	}
	fmt.Fprint(w, "success new post")
}
func handlerPostDelete(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Path[len("/admin/post/"):]
	if !verifyUUID(uuid) {
		http.Error(w, "uuid not verified", 400)
		return
	}
	err := postDelete(uuid)
	if err != nil {
		log.Println(err)
		http.Error(w, "fail to delete post", 500)
		return
	}
	fmt.Fprint(w, "success delete post")
}
func handlerPostUpdate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uuid := r.Form.Get("uuid")
	if !verifyUUID(uuid) {
		http.Error(w, "uuid not verified", 400)
		return
	}
	title := r.Form.Get("title")
	date := r.Form.Get("date")
	content := r.Form.Get("content")
	brief := r.Form.Get("brief")
	requestTagUUIDs := r.Form.Get("tags")
	categoryUUID := r.Form.Get("category")
	tagUUIDs, err := parseStrings(requestTagUUIDs)
	if err != nil {
		log.Println(err)
		http.Error(w, "fail to parse tags", 500)
		return
	}
	err = postUpdate(uuid, title, date, brief, content, tagUUIDs, categoryUUID)
	if err != nil {
		log.Println(err)
		http.Error(w, "fail to new post", 500)
		return
	}
	fmt.Fprint(w, "success update post")
}
func handlerTagNew(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	name := r.Form.Get("name")
	if name == "" {
		http.Error(w, "lack name", 400)
		return nil
	}
	err := tagNew(name)
	if err != nil {
		http.Error(w, "fail to new tag", 500)
		return err
	}
	return nil
}
func handlerTagUpdate(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	newName := r.Form.Get("name")
	uuid := r.Form.Get("uuid")
	if !verifyUUID(uuid) {
		http.Error(w, "uuid not verified", 400)
		return nil
	}
	if newName == "" {
		http.Error(w, "new name not found", 400)
		return nil
	}
	err := tagUpdate(uuid, newName)
	if err != nil {
		return err
	}
	return nil
}
func handlerTagDelete(w http.ResponseWriter, r *http.Request) error {
	uuid := r.URL.Path[len("/admin/tag"):]
	forceString := r.URL.Query().Get("force")
	var force bool
	if forceString == "true" {
		force = true
	} else {
		force = false
	}
	if !verifyUUID(uuid) {
		http.Error(w, "uuid not verified", 400)
		return nil
	}
	err := tagDelete(uuid, force)
	if err != nil {
		return err
	}
	return nil
}
