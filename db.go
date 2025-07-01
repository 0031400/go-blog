package main

import (
	"database/sql"
	"errors"
	"fmt"

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
	newUUID, err := getNewUUID("posts")
	if err != nil {
		return err
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
	err = db.QueryRow("SELECT title,date,brief,content,categoryUUID FROM posts WHERE uuid = ?", uuid).Scan(&newTitle, &newDate, &newBrief, &newContent, &newCategoryUUID)
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
	if tagUUIDs == nil {
		return nil
	}
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
	_, err = tx.Exec("UPDATE posts SET title = ?, date = ?, brief = ?, content = ?, categoryUUID = ? WHERE uuid = ?", newTitle, newDate, newBrief, newContent, newCategoryUUID, uuid)
	if err != nil {
		return err
	}
	return nil
}
func tagNew(name string) error {
	newUUID, err := getNewUUID("tags")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO tags (name,uuid) VALUES (?,?)", name, newUUID)
	if err != nil {
		return err
	}
	return nil
}
func getNewUUID(table string) (string, error) {
	var newUUID string
	var err error
	newUUIDTime := 0
	for {
		var dummy int
		newUUID, err = uuid()
		if err != nil {
			return "", err
		}
		query := fmt.Sprintf("SELECT 1 FROM %s WHERE uuid = ? LIMIT 1", table)
		err = db.QueryRow(query, newUUID).Scan(&dummy)
		if err == sql.ErrNoRows {
			return newUUID, nil
		} else if err != nil {
			return "", err
		}
		newUUIDTime++
		if newUUIDTime == 6 {
			return "", errors.New("fail to generate new uuid")
		}
	}
}
func tagUpdate(uuid string, newName string) error {
	err := uuidOnlyOne(uuid, "tags")
	if err != nil {
		return err
	}
	rowsAffected, err := db.Exec("UPDATE tags SET name = ? WHERE uuid = ?", newName, uuid)
	if err != nil {
		return err
	}
	affectedNum, err := rowsAffected.RowsAffected()
	if err != nil {
		return err
	}
	if affectedNum != 1 {
		return fmt.Errorf("rows affected num not right.num: %d", affectedNum)
	}
	return nil
}
func tagDelete(uuid string, force bool) error {
	err := uuidOnlyOne(uuid, "tag")
	if err != nil {
		return err
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
	var num int
	err = db.QueryRow("SELECT COUNT (*) FROM post_tags WHERE tagUUID = ?", uuid).Scan(&num)
	if err != nil {
		return err
	}
	if num != 0 {
		if force == false {
			return errors.New("exist posts attached with the tag")
		} else {
			_, err = tx.Exec("DELETE FROM post_tags WHERE tagUUID = ?", uuid)
			if err != nil {
				return err
			}
		}
	}
	_, err = tx.Exec("DELETE FROM tags WHERE uuid = ?", uuid)
	if err != nil {
		return err
	}
	return nil
}
func uuidOnlyOne(uuid string, table string) error {
	var num int
	err := db.QueryRow(fmt.Sprintf("SELECT 1 FROM %s WHERE uuid = ?", table), uuid).Scan(&num)
	if err != nil {
		return err
	}
	if num != 1 {
		return fmt.Errorf("tag num not right.num: %d", num)
	}
	return nil
}
