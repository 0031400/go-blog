package main

type post struct {
	uuid         string
	title        string
	date         string
	brief        string
	content      string
	tagUUIDs     []string
	categoryUUID string
}
