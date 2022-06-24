package main

import (
	"encoding/hex"
	"strconv"
	"strings"
)

func DecToHexstring(hex_num uint64) string {
	return strconv.FormatInt(int64(hex_num), 16)
}

func BinstringToString(byteArr string) string {
	byteArr = strings.ReplaceAll(byteArr, ",", "")
	byteArr = strings.ReplaceAll(byteArr, "00", "")

	data, err := hex.DecodeString(byteArr)
	if err != nil {
		panic(err)
	}

	return string(data)
}
