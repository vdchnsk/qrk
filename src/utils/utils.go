package utils

import (
	"encoding/binary"
	"fmt"
	"sort"
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

func ReadUint8(bytes []byte) uint8 {
	return uint8(bytes[0])
}

func ReadUint16(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}

func ReadUint32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

func SortByString[T fmt.Stringer](list []T) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].String() < list[j].String()
	})
}
