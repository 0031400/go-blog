package main

type Category struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}
type Tag struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}
type Post struct {
	UUID        string   `json:"uuid"`
	Title       string   `json:"name"`
	Date        string   `json:"date"`
	Brief       string   `json:"brief"`
	Content     string   `json:"content"`
	TheCategory Category `json:"category"`
	TheTags     []Tag    `json:"tags"`
}
