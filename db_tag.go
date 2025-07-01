package main

import "errors"

func tagUpdate(uuid string, newName string) error {
	_, err := db.Exec("UPDATE tags SET name = ? WHERE uuid = ?", newName, uuid)
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
func tagDelete(uuid string, force bool) error {
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
	err = db.QueryRow("SELECT COUNT (*) FROM post_tags WHERE tagUUID = ? AND deleteAt IS NULL", uuid).Scan(&num)
	if err != nil {
		return err
	}
	if num != 0 {
		if force == false {
			return errors.New("exist posts attached with the tag")
		} else {
			_, err = tx.Exec("UPDATE post_tags SET deleteAt = ? WHERE tagUUID = ?", nowDate(), uuid)
			if err != nil {
				return err
			}
		}
	}
	_, err = tx.Exec("UPDATE tags SET deleteAt = ? WHERE uuid = ?", nowDate(), uuid)
	if err != nil {
		return err
	}
	return nil
}
