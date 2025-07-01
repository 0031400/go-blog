package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"
)

func parseStrings(s string) ([]string, error) {
	if s == "" {
		return nil, nil
	}
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
func stringInclude(s string, list []string) int {
	for i, v := range list {
		if v == s {
			return i
		}
	}
	return -1
}
func stringUpdateIfNotNull(oldS *string, newS string) {
	if newS != "" {
		log.Println("try to change: ", oldS, newS)
		*oldS = newS
	}
}
func nowDate() string {
	return time.Now().Format("20060102")
}
