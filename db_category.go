package main

import "errors"

func categoryNew(name string) error {
	newUUID, err := getNewUUID("categories")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO categories (name,uuid) VALUES (?,?)", name, newUUID)
	if err != nil {
		return err
	}
	return nil
}
func categoryDelete(uuid string, force bool) error {
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
	err = db.QueryRow("SELECT COUNT (*) FROM posts WHERE categoryUUID = ? AND deleteAt IS NOT NULL", uuid).Scan(&num)
	if err != nil {
		return err
	}
	if num != 0 && force == false {
		return errors.New("exist posts attached with the category")
	}
	_, err = tx.Exec("UPDATE categories SET deleteAt = ? WHERE uuid = ?", nowDate(), uuid)
	if err != nil {
		return err
	}
	return nil
}
func categoryUpdate(uuid string, newName string) error {
	_, err := db.Exec("UPDATE categories SET name = ? WHERE uuid = ?", newName, uuid)
	if err != nil {
		return err
	}
	return nil
}
