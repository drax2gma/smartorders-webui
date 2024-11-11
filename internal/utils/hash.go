package utils

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
)

func GenerateXXH64Hash(input string) string {
	hash := xxhash.Sum64String(input)
	return fmt.Sprintf("%016x", hash)
}
