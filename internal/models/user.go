package models

import (
	"encoding/binary"
	"fmt"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	prime1 uint32 = 2654435761
	prime2 uint32 = 2246822519
	prime3 uint32 = 3266489917
	prime4 uint32 = 668265263
	prime5 uint32 = 374761393
)

func rotl32(x, r uint32) uint32 {
	return (x << r) | (x >> (32 - r))
}

func GenerateUserID(email string) string {
	var h uint32 = prime5
	data := []byte(email)
	len := len(data)

	for len >= 4 {
		k := binary.LittleEndian.Uint32(data)
		h += k * prime3
		h = rotl32(h, 17) * prime4
		data = data[4:]
		len -= 4
	}

	for len > 0 {
		h += uint32(data[0]) * prime5
		h = rotl32(h, 11) * prime1
		data = data[1:]
		len--
	}

	h ^= h >> 15
	h *= prime2
	h ^= h >> 13
	h *= prime3
	h ^= h >> 16

	return fmt.Sprintf("%08x", h)
}
