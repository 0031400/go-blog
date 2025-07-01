package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func handlerPostNew(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	title := r.Form.Get("title")
	date := r.Form.Get("date")
	content := r.Form.Get("content")
	brief := r.Form.Get("brief")
	requestTagUUIDs := r.Form.Get("tags")
	categoryUUID := r.Form.Get("category")
	tagUUIDs, err := parseStrings(requestTagUUIDs)
	if err != nil {
		http.Error(w, "fail to parse tags", 500)
		return err
	}
	if !verifyDate(date) {
		http.Error(w, "date not verified", 400)
		return nil
	}
	err = postNew(title, date, brief, content, tagUUIDs, categoryUUID)
	if err != nil {
		http.Error(w, "fail to new post", 500)
		return err
	}
	fmt.Fprint(w, "success new post")
	return nil
}
func handlerPostDelete(w http.ResponseWriter, r *http.Request) error {
	uuid := r.URL.Query().Get("uuid")
	if !verifyUUID(uuid) {
		http.Error(w, "uuid not verified", 400)
		return nil
	}
	err := postDelete(uuid)
	if err != nil {
		http.Error(w, "fail to delete post", 500)
		return err
	}
	fmt.Fprint(w, "success delete post")
	return nil
}
func handlerPostUpdate(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	uuid := r.Form.Get("uuid")
	if !verifyUUID(uuid) {
		http.Error(w, "uuid not verified", 400)
		return nil
	}
	title := r.Form.Get("title")
	date := r.Form.Get("date")
	content := r.Form.Get("content")
	brief := r.Form.Get("brief")
	requestTagUUIDs := r.Form.Get("tags")
	categoryUUID := r.Form.Get("category")
	tagUUIDs, err := parseStrings(requestTagUUIDs)
	if err != nil {
		http.Error(w, "fail to parse tags", 500)
		return err
	}
	err = postUpdate(uuid, title, date, brief, content, tagUUIDs, categoryUUID)
	if err != nil {
		http.Error(w, "fail to new post", 500)
		return err
	}
	fmt.Fprint(w, "success update post")
	return nil
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
	fmt.Fprint(w, "success to new tag")
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
		fmt.Fprint(w, "fail to update tag")
		return err
	}
	fmt.Fprint(w, "success to update tag")
	return nil
}
func handlerTagDelete(w http.ResponseWriter, r *http.Request) error {
	uuid := r.URL.Query().Get("uuid")
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
		http.Error(w, "fail to delete tag", 500)
		return err
	}
	fmt.Fprint(w, "success to delete tag")
	return nil
}
func handlerCategoryNew(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	name := r.Form.Get("name")
	if name == "" {
		http.Error(w, "lack name", 400)
		return nil
	}
	err := categoryNew(name)
	if err != nil {
		http.Error(w, "fail to new category", 500)
		return err
	}
	fmt.Fprint(w, "success to new category")
	return nil
}
func handlerCategoryUpdate(w http.ResponseWriter, r *http.Request) error {
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
	err := categoryUpdate(uuid, newName)
	if err != nil {
		fmt.Fprint(w, "fail to update category")
		return err
	}
	fmt.Fprint(w, "success to update category")
	return nil
}
func handlerCategoryDelete(w http.ResponseWriter, r *http.Request) error {
	uuid := r.URL.Query().Get("uuid")
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
	err := categoryDelete(uuid, force)
	if err != nil {
		http.Error(w, "fail to delete category", 500)
		return err
	}
	fmt.Fprint(w, "success to delete category")
	return nil

}
func handlerCategoryList(w http.ResponseWriter, r *http.Request) error {
	uuid := r.URL.Query().Get("uuid")
	if !verifyUUID(uuid) {
		http.Error(w, "uuid not verified", 400)
		return nil
	}
	index, err := strconv.Atoi(r.URL.Query().Get("index"))
	if err != nil {
		http.Error(w, "fail to parse index", 500)
		return err
	}
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		http.Error(w, "fail to parse size", 500)
		return err
	}
	if index <= 0 || size <= 0 {
		http.Error(w, "query variable not allowed", 400)
		return nil
	}
	postList, err := categoryList(uuid, index, size)
	if err != nil {
		http.Error(w, "fail to get post list", 500)
		return err
	}
	err = jsonResponse(w, postList)
	if err != nil {
		http.Error(w, "fail to response json", 500)
	}
	return nil
}

func handlerTagList(w http.ResponseWriter, r *http.Request) error {
	uuid := r.URL.Query().Get("uuid")
	if !verifyUUID(uuid) {
		http.Error(w, "uuid not verified", 400)
		return nil
	}
	index, err := strconv.Atoi(r.URL.Query().Get("index"))
	if err != nil {
		http.Error(w, "fail to parse index", 500)
		return err
	}
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		http.Error(w, "fail to parse size", 500)
		return err
	}
	if index <= 0 || size <= 0 {
		http.Error(w, "query variable not allowed", 400)
		return nil
	}
	postList, err := tagList(uuid, index, size)
	if err != nil {
		http.Error(w, "fail to get post list", 500)
		return err
	}
	err = jsonResponse(w, postList)
	if err != nil {
		http.Error(w, "fail to response json", 500)
	}
	return nil
}

func handlerPostList(w http.ResponseWriter, r *http.Request) error {
	index, err := strconv.Atoi(r.URL.Query().Get("index"))
	if err != nil {
		http.Error(w, "fail to parse index", 500)
		return err
	}
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		http.Error(w, "fail to parse size", 500)
		return err
	}
	if index <= 0 || size <= 0 {
		http.Error(w, "query variable not allowed", 400)
		return nil
	}
	postList, err := postList(index, size)
	if err != nil {
		http.Error(w, "fail to get post list", 500)
		return err
	}
	err = jsonResponse(w, postList)
	if err != nil {
		http.Error(w, "fail to response json", 500)
	}
	return nil
}

func handlerPostDetail(w http.ResponseWriter, r *http.Request) error {
	uuid := r.URL.Query().Get("uuid")
	if !verifyUUID(uuid) {
		http.Error(w, "uuid not verified", 400)
		return nil
	}
	thePostDetail, err := postDetail(uuid)
	if err != nil {
		http.Error(w, "fail to get post detail", 400)
		return err
	}
	err = jsonResponse(w, thePostDetail)
	if err != nil {
		http.Error(w, "fail to response json", 500)
		return err
	}
	return nil
}
func jsonResponse(w http.ResponseWriter, item any) error {
	jsonData, err := json.Marshal(item)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonData))
	return nil
}
