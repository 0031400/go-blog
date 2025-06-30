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
	var newTitle, newBrief, newDate, newContent, newCategoryUUID string
	err = db.QueryRow("SELECT title,date,brief,content,categoryUUID FROM posts WHERE uuid =?", uuid).Scan(&newTitle, &newDate, &newBrief, &newContent, &newCategoryUUID)
	if err == sql.ErrNoRows {
		return errors.New("post no founded")
	} else if err != nil {
		return err
	}
	stringUpdateIfNotNull(&newTitle, title)
	stringUpdateIfNotNull(&newDate, date)
	stringUpdateIfNotNull(&newBrief, brief)
	stringUpdateIfNotNull(&newContent, content)
	stringUpdateIfNotNull(&newCategoryUUID, categoryUUID)
	// del with the tags
	flagList := make([]bool, len(tagUUIDs))
	rows, err := tx.Query("SELECT tagUUID FROM post_tags WHERE postUUID = ?", uuid)
	if err != nil {
		return err
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		var tagUUID string
		rows.Scan(&tagUUID)
		index := stringInclude(tagUUID, tagUUIDs)
		if index != -1 {
			// if exist
			flagList[i] = true
		} else {
			// delete the flag
			_, err = tx.Exec("DELETE FROM post_tags WHERE postUUID = ? AND tagUUID =?", uuid, tagUUID)
			if err != nil {
				return err
			}
		}
		i++
	}
	for i, v := range tagUUIDs {
		if flagList[i] {
			continue
		}
		_, err := tx.Exec("INSERT INTO post_tags (postUUID,tagUUId) VALUES (?,?)", uuid, v)
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec("UPDATE posts SET title = ?, date = ?, brief = ?, content = ?, categoryUUID = ? WHERE uuid = ?", title, date, brief, content, categoryUUID, uuid)
	if err != nil {
		return err
	}
	return nil
}
