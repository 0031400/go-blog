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

func postNew(title string, date string, brief string, content string, tagUUIDs []string, categoryUUID string) error {
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
	newUUID, err := getNewUUID("posts")
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO posts(uuid,title,brief,content,date,categoryUUID) VALUES (?,?,?,?,?,?)", newUUID, title, brief, content, date, categoryUUID)
	if err != nil {
		return err
	}
	for _, tagUUID := range tagUUIDs {
		_, err = tx.Exec("INSERT INTO post_tags(postUUID,tagUUID) VALUES(?,?)", newUUID, tagUUID)
		if err != nil {
			return err
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
	_, err = tx.Exec("UPDATE post_tags SET deleteAt = ? WHERE postUUID = ?", nowDate(), uuid)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE posts SET deleteAt = ? WHERE uuid = ?", nowDate(), uuid)
	if err != nil {
		return err
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
	err = db.QueryRow("SELECT title,date,brief,content,categoryUUID FROM posts WHERE uuid = ? AND deleteAt IS NOT NULL", uuid).Scan(&newTitle, &newDate, &newBrief, &newContent, &newCategoryUUID)
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
	// deal with the tags
	flagList := make([]bool, len(tagUUIDs))
	rows, err := tx.Query("SELECT tagUUID FROM post_tags WHERE postUUID = ? AND deleteAt IS NOT NULL", uuid)
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
			_, err = tx.Exec("UPDATE post_tags SET deleteAt = ? WHERE postUUID = ? AND tagUUID =?", nowDate(), uuid, tagUUID)
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

func postList(index int, size int) ([]post, error) {
	var uuid, title, date, brief, categoryUUID, categoryName, tagUUID, tagName string
	var thePostList []post
	var theTagList []tag
	rows, err := db.Query("SELECT uuid,title,date,brief,categoryUUID FROM posts ORDER BY date DESC LIMIT ? OFFSET ?", (index-1)*size, size)
	if err != nil {
		return []post{}, err
	}
	defer rows.Close()
	for rows.Next() {
		theTagList = theTagList[:0]
		err = rows.Scan(&uuid, &title, &date, &brief, &categoryUUID)
		if err != nil {
			return []post{}, err
		}
		err = db.QueryRow("SELECT name FROM categories WHERE uuid = ?", categoryUUID).Scan(&categoryName)
		if err != nil {
			return []post{}, err
		}
		tagRows, err := db.Query("SELECT tagUUID FROM post_tags WHERE postUUID = ?", uuid)
		if err != nil {
			return []post{}, err
		}
		defer tagRows.Close()
		for tagRows.Next() {
			err = tagRows.Scan(&tagUUID)
			if err != nil {
				return []post{}, err
			}
			err = db.QueryRow("SELECT name FROM tags WHERE uuid = ?", tagUUID).Scan(&tagName)
			if err != nil {
				return []post{}, err
			}
			theTagList = append(theTagList, tag{name: tagName, uuid: tagUUID})
		}
		thePostList = append(thePostList, post{title: title, date: date, uuid: uuid, brief: brief, theTags: theTagList, theCategory: category{name: categoryName, uuid: categoryUUID}})
	}
	return thePostList, nil
}
