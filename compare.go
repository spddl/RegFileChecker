package main

// CompareInteger
func CompareDWord(regValue, fileValue string) bool {
	switch {
	case regValue == "0" && fileValue == "":
		fallthrough
	case regValue == fileValue:
		return true
	default:
		return false
	}
}
