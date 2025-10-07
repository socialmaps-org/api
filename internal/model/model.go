package model

import (
	"database/sql"
	"math/rand"
	"strings"
)

type Model struct {
	DB *sql.DB
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const n = 32

func randomID(prefix string) string {
	var b strings.Builder
	b.Grow(len(prefix) + n)

	_, err := b.WriteString(prefix)
	if err != nil {
		panic(err)
	}

	for i := 0; i < n; i++ {
		err = b.WriteByte(alphabet[rand.Int63n(int64(len(alphabet)))])
		if err != nil {
			panic(err)
		}
	}

	return b.String()
}
