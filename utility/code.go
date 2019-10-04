package utility

import (
	"math/rand"
	"strings"
	"time"
)

// GenerateCode is generate random numeric code
func GenerateCode(length int, characters string) string {
	rand.Seed(time.Now().UnixNano())
	var code []string
	for index := 0; index < length; index++ {
		code = append(code, string(characters[rand.Intn(len(characters))]))
	}
	return strings.Join(code, "")
}
