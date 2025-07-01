package main

import (
	"database/sql"
	"errors"
	"fmt"
)

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
