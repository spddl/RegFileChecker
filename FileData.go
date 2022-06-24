package main

type FileData struct {
	Value string
	Key   string
	Type  int
}

func (fd *FileData) Clear() {
	fd = &FileData{}
}

// const (
// 	SZ = iota
// 	HEX
// 	DWORD
// 	HEX0
// 	HEX1
// 	HEX2
// 	HEX3
// 	HEX4
// 	HEX5
// 	HEX7
// 	HEX8
// 	HEXA
// 	HEXB
// )
