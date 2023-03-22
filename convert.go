package main

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"strings"
)

// uint64 => string
func DecToHexstring(hex_num uint64) string {
	return strconv.FormatInt(int64(hex_num), 16)
}

// string => string
func BinstringToString(byteArr string) string {
	byteArr = strings.ReplaceAll(byteArr, ",", "")
	byteArr = strings.ReplaceAll(byteArr, "00", "")

	data, err := hex.DecodeString(byteArr)
	if err != nil {
		panic(err)
	}

	return string(data)
}

// string => []byte
func BinStringToByteArray(byteArr string) []byte {
	byteArr = strings.ReplaceAll(byteArr, ",", "")

	data, err := hex.DecodeString(byteArr)
	if err != nil {
		panic(err)
	}

	return data
}

// []byte => string
func ByteArrayToBinString(byteArr []byte) string {
	data := hex.EncodeToString(byteArr)

	var result string
	var i = 0
	for _, char := range data {
		if i%2 == 0 {
			result += ","
		}
		result += string(char)
		i++
	}
	if result[0] == 44 { // Comma found
		return result[1:]
	}
	return result
}

// []byte => uint64
func ByteArrayToUint64(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}

// uint64 => []byte
func Uint64ToByteArray(data uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, data)
	return b
}

// CompareInteger
func CompareDWord(regValue, fileValue string) bool {
	switch {
	case (regValue == "0" && fileValue == "") || regValue == fileValue:
		return true
	default:
		return false
	}
}
