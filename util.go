package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
)

func parseStrings(s string) ([]string, error) {
	var result []string
	err := json.Unmarshal([]byte(s), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func verifyDate(d string) bool {
	return len(d) == 8
}
func verifyUUID(u string) bool {
	return len(u) == 32
}
func uuid() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(u), nil
}
