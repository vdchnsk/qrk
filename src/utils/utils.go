package utils

import (
	"encoding/binary"
	"strings"
)

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func RemoveWhitespaces(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

func MaxIntForBytes(amountOfBytes int) int {
	return (1 << (amountOfBytes * 8)) - 1
}

func ReadUint16(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}

func ReadUint32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}
