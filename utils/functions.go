package utils

import (
	"github.com/google/uuid"
	"strings"
)

// NewUUID 生成 UUID
func NewUUID() string {
	uid := uuid.Must(uuid.NewRandom())
	var idBytes [32]byte
	copy(idBytes[:], strings.Replace(uid.String(), "-", "", -1))
	return string(idBytes[:])
}
