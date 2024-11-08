package common

import (
	"fmt"
	"strings"
)

func InitTranslator() {
	// TODO
}

func Translate(key string, args ...interface{}) string {
	return fmt.Sprintf(strings.Split(key, ":")[1], args...)
}
