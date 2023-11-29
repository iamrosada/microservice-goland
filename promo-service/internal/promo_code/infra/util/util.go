package util

import "time"

func GenerateNewID() uint {
	return uint(time.Now().UnixNano())
}
