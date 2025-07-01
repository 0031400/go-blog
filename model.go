package main

type category struct {
	name string
	uuid string
}
type tag struct {
	name string
	uuid string
}
type post struct {
	uuid        string
	title       string
	date        string
	brief       string
	theCategory category
	theTags     []tag
}
