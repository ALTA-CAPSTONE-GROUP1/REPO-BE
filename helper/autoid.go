package helper

import "fmt"

func GenerateIDFromPositionTag(positionTag string) string {
	var lastID int

	lastID++

	id := fmt.Sprintf("%s%02d", positionTag, lastID)

	return id
}
