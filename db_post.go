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
	err = db.QueryRow("SELECT title,date,brief,content,categoryUUID FROM posts WHERE uuid = ? AND deleteAt IS NULL", uuid).Scan(&newTitle, &newDate, &newBrief, &newContent, &newCategoryUUID)
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
	rows, err := tx.Query("SELECT tagUUID FROM post_tags WHERE postUUID = ? AND deleteAt IS NULL", uuid)
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
func postListQuery(detail bool, query string, args ...any) ([]Post, error) {
	var uuid, title, date, brief, categoryUUID, content string
	var thePostList []Post
	var theTagList []Tag
	rows, err := db.Query(query, args...)
	if err != nil {
		return []Post{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&uuid, &title, &date, &brief, &content, &categoryUUID)
		if err != nil {
			return []Post{}, err
		}
		oneCategory, err := fromUUIDToCategory(categoryUUID)
		if err != nil {
			return []Post{}, err
		}
		theTagList, err = fromPostUUIDToTags(uuid)
		if err != nil {
			return []Post{}, err
		}
		thePost := Post{Title: title, Date: date, UUID: uuid, Brief: brief, Content: content, TheTags: theTagList, TheCategory: oneCategory}
		if detail {
			return []Post{thePost}, nil
		} else {
			thePost.Content = ""
		}
		thePostList = append(thePostList, thePost)
	}
	return thePostList, nil
}
func postList(index int, size int) ([]Post, error) {
	return postListQuery(false, "SELECT uuid,title,date,brief,content,categoryUUID FROM posts WHERE deleteAt IS NULL ORDER BY date DESC LIMIT ? OFFSET ?", size, (index-1)*size)
}
func categoryList(uuid string, index int, size int) ([]Post, error) {
	return postListQuery(false, "SELECT uuid,title,date,brief,content,categoryUUID FROM posts WHERE categoryUUID = ? AND deleteAt IS NULL ORDER BY date DESC LIMIT ? OFFSET ?", uuid, size, (index-1)*size)
}
func postDetail(uuid string) (Post, error) {
	thePostList, err := postListQuery(true, "SELECT uuid,title,date,brief,content,categoryUUID FROM posts WHERE uuid = ? AND deleteAt IS NULL", uuid)
	if err != nil {
		return Post{}, err
	}
	return thePostList[0], nil
}
func tagList(tagUUID string, index int, size int) ([]Post, error) {
	var thePostList []Post
	var postUUID string
	tagRows, err := db.Query("SELECT postUUID FROM post_tags WHERE tagUUID = ? AND deleteAt IS NULL LIMIT ? OFFSET ?", tagUUID, size, size*(index-1))
	if err != nil {
		return []Post{}, err
	}
	defer tagRows.Close()
	for tagRows.Next() {
		err = tagRows.Scan(&postUUID)
		if err != nil {
			return []Post{}, err
		}
		onePost, err := fromUUIDToPost(postUUID)
		if err != nil {
			return []Post{}, err
		}
		thePostList = append(thePostList, onePost)
	}
	return thePostList, nil
}
func categoryFromUUIDToName(uuid string) (string, error) {
	var name string
	err := db.QueryRow("SELECT name FROM categories WHERE uuid = ?", uuid).Scan(&name)
	return name, err
}

func fromUUIDToPost(postUUID string) (Post, error) {
	onePostList, err := postListQuery(false, "SELECT uuid,title,date,brief,content,categoryUUID FROM posts WHERE deleteAt IS NULL AND uuid = ?", postUUID)
	if err != nil {
		return Post{}, err
	}
	return onePostList[0], nil
}
func TagFromUUIDToName(uuid string) (string, error) {
	var name string
	err := db.QueryRow("SELECT name FROM tags WHERE uuid = ?", uuid).Scan(&name)
	return name, err
}
func fromUUIDToCategory(uuid string) (Category, error) {
	name, err := categoryFromUUIDToName(uuid)
	if err != nil {
		return Category{}, err
	}
	return Category{UUID: uuid, Name: name}, nil
}
func fromPostUUIDToTags(postUUID string) ([]Tag, error) {
	var theTagList []Tag
	var tagUUID, tagName string
	tagRows, err := db.Query("SELECT tagUUID FROM post_tags WHERE postUUID = ? AND deleteAt IS NULL", postUUID)
	if err != nil {
		return []Tag{}, err
	}
	defer tagRows.Close()
	for tagRows.Next() {
		err = tagRows.Scan(&tagUUID)
		if err != nil {
			return []Tag{}, err
		}
		err = db.QueryRow("SELECT name FROM tags WHERE uuid = ?", tagUUID).Scan(&tagName)
		if err != nil {
			return []Tag{}, err
		}
		theTagList = append(theTagList, Tag{Name: tagName, UUID: tagUUID})
	}
	return theTagList, nil
}
