package main

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	theDB, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		panic(err)
	}
	db = theDB
}
func tagExist(tagName string) (bool, error) {
	query := "SELECT COUNT(*) FROM tags WHERE uuid =?"
	var count int
	err := db.QueryRow(query, tagName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func categoryExist(tagName string) (bool, error) {
	query := "SELECT COUNT(*) FROM categories WHERE uuid =?"
	var count int
	err := db.QueryRow(query, tagName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func tagsExist(tagNames []string) (bool, error) {
	for _, v := range tagNames {
		exist, err := tagExist(v)
		if err != nil {
			return false, err
		}
		if exist == false {
			return false, nil
		}
	}
	return true, nil
}
func postNew(title string, date string, brief string, content string, tagUUIDs []string, categoryUUID string) error {
	var newUUID string
	var err error
	newUUIDTime := 0
	for {
		var dummy int
		newUUID, err = uuid()
		if err != nil {
			return err
		}
		query := "SELECT 1 FROM posts WHERE uuid = ? LIMIT 1"
		err = db.QueryRow(query, newUUID).Scan(&dummy)
		if err == sql.ErrNoRows {
			break
		} else if err != nil {
			return err
		}
		newUUIDTime++
		if newUUIDTime == 6 {
			return errors.New("fail to generate new uuid")
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	result, err := tx.Exec("INSERT INTO posts(uuid,title,brief,content,date,categoryUUID) VALUES (?,?,?,?,?,?)", newUUID, title, brief, content, date, categoryUUID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected != 1 {
		return errors.New("fail to insert posts")
	}
	for _, tagUUID := range tagUUIDs {
		result, err := tx.Exec("INSERT INTO post_tags(postUUID,tagUUID) VALUES(?,?)", newUUID, tagUUID)
		if err != nil {
			return err
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected != 1 {
			return errors.New("fail to insert post_tags")
		}
	}
	return nil
}
func postDelete(uuid string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	_, err = tx.Exec("DELETE FROM post_tags WHERE postUUID = ?", uuid)
	if err != nil {
		return err
	}
	result, err := tx.Exec("DELETE FROM posts WHERE uuid = ?", uuid)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected != 1 {
		return errors.New("fail to delete posts")
	}
	return nil
}
func postUpdate(uuid string, title string, date string, brief string, content string, tagUUIDs []string, categoryUUID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	var oldPost post
	err = db.QueryRow("SELECT * FROM posts WHERE uuid =?", uuid).Scan(&oldPost)
	if err == sql.ErrNoRows {
		return errors.New("post no founded")
	} else if err != nil {
		return err
	}
	newPost := post{
		uuid:         oldPost.uuid,
		title:        oldPost.title,
		date:         oldPost.date,
		brief:        oldPost.brief,
		content:      oldPost.content,
		categoryUUID: oldPost.categoryUUID,
	}
	rows, err := db.Query("SELECT tagUUID FROM post_tags WHERE postUUID = ?", uuid)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var tagUUID string
		rows.Scan(&tagUUID)
		exist := false
		for _, i := range tagUUIDs {
			if i == tagUUID {
				exist = true
				break
			}
		}
		if !exist {
			tx.Exec("DELETE FROM post_tags WHERE postUUID = ? AND tagUUID =?", uuid, tagUUID)
		}
	}
}
