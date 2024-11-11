package utils

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
)

const (
	masterSalt = "be0d748d-58e2-50d5-b194-8996858d97ba-7d30b755-e3ec-5f0d-a2e8-0a003198c34d"
)

func GenerateXXH64Hash(input string) string {
	hash := xxhash.Sum64String(input + masterSalt)
	return fmt.Sprintf("%016x", hash)
}
