package model

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"math"
	"strings"
	"time"

	"codeberg.org/socialmaps/api/internal/mytime"
)

const (
	ver            = "1"
	sep            = "_"
	base62alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// ⌈log₆₂ 2¹²⁸⌉ = 22 chars
	randLen = 22
)

type ID struct {
	kind    string
	version string
	time    time.Time
	rand    string
}

var (
	// We use base32hex because it preserves the bitwise sort order.
	timeEnc = base32.HexEncoding.WithPadding(base32.NoPadding)
	timeEnd = binary.BigEndian
)

func NewRandomID(kind string) *ID {
	return &ID{
		kind:    kind,
		version: ver,
		time:    mytime.Now(),
		rand:    randText(),
	}
}

func EarliestID(kind string) *ID {
	return &ID{
		kind:    kind,
		version: ver,
		time:    time.Unix(0, 0),
		rand:    strings.Repeat("0", randLen),
	}
}

func LatestID(kind string) *ID {
	return &ID{
		kind:    kind,
		version: ver,
		time:    time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
		rand:    strings.Repeat("z", randLen),
	}
}

func (id *ID) String() string {
	var bld strings.Builder

	bld.Grow(len(id.kind) + len(sep) + len(ver) + timeEnc.EncodedLen(8) + 22)

	bld.WriteString(id.kind)
	bld.WriteString(sep)
	bld.WriteString(ver)

	var timeBytes [8]byte
	timeEnd.PutUint64(timeBytes[:], uint64(id.time.Unix()))
	bld.WriteString(timeEnc.EncodeToString(timeBytes[:]))

	bld.WriteString(id.rand)

	return bld.String()
}

func ParseID(id string) *ID {
	tokens := strings.Split(id, sep)
	if len(tokens) != 2 {
		panic("invalid format")
	}

	kind, tail := tokens[0], tokens[1]

	ver, body := tail[:1], tail[1:]

	timeLen := timeEnc.EncodedLen(8)
	timeString, random := body[:timeLen], body[timeLen:]

	timeBytes, err := timeEnc.DecodeString(timeString)
	if err != nil {
		panic(err)
	}
	sec := timeEnd.Uint64(timeBytes)
	if sec > math.MaxInt64 {
		panic("unsafe cast")
	}
	ts := time.Unix(int64(sec), 0).UTC()

	return &ID{
		kind:    kind,
		version: ver,
		time:    ts,
		rand:    random,
	}
}

// randText returns a cryptographically random string using the base62 alphabet
// (like base64 but containing only alphanumeric characters).
// The result contains at least 128 bits of randomness, enough to make the
// likelihood of collisions vanishingly small.
// Copied and adapted from Go Standard library crypto/rand#Text()
func randText() string {
	// ⌈log₆₂ 2¹²⁸⌉ = 22 chars
	src := make([]byte, 22)
	rand.Read(src)
	for i := range src {
		src[i] = base62alphabet[src[i]%62]
	}
	return string(src)
}
