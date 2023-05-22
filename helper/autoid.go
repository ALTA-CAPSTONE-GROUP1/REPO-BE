package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GenerateIDFromPositionTag(positionTag string) string {
	var lastID int

	lastID++

	id := fmt.Sprintf("%s%02d", positionTag, lastID)

	return id
}

func GenerateUniqueSign(userID string) (string, error) {
	const (
		letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		signLength  = 5
	)

	rand.Seed(time.Now().UnixNano())
	randomStr := make([]byte, signLength)
	for i := range randomStr {
		randomStr[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	sign := userID + strings.ToUpper(string(randomStr))
	return sign, nil
}
