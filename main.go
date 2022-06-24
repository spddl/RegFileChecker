package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"golang.org/x/sys/windows/registry"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	file, err := os.Open(RegFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// scanner := bufio.NewScanner(transform.NewReader(file, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()))
	scanner := bufio.NewScanner(file)
	var WorkingRegPath = new(Data)

	var index uint

	var line string
	for scanner.Scan() {
		index++
		if index == 1 {
			continue // skip "Windows Registry Editor Version 5.00"
		}

		if rawLine := scanner.Text(); len(rawLine) > 1 {
			// fmt.Printf("[%d] %s\n", index, rawLine)
			if rawLine[:1] == ";" {
				continue // skip comments
			}

			if strings.HasSuffix(rawLine, "\\") {
				line += strings.TrimSuffix(strings.TrimLeftFunc(rawLine, trimSpace), "\\")
				continue
			} else {
				line += strings.TrimLeftFunc(rawLine, trimSpace)
			}

			if line[:1] == "[" {
				WorkingRegPath = new(Data)
				line = line[1 : len(line)-1]

				var rootkeySplit = strings.SplitN(line, `\`, 2)
				WorkingRegPath.StringPath = rootkeySplit[1]

				switch rootkeySplit[0] {
				case "HKEY_CLASSES_ROOT":
					WorkingRegPath.RootKey = registry.CLASSES_ROOT
				case "HKEY_CURRENT_USER":
					WorkingRegPath.RootKey = registry.CURRENT_USER
				case "HKEY_LOCAL_MACHINE":
					WorkingRegPath.RootKey = registry.LOCAL_MACHINE
				case "HKEY_USERS":
					WorkingRegPath.RootKey = registry.USERS
				case "HKEY_CURRENT_CONFIG":
					WorkingRegPath.RootKey = registry.CURRENT_CONFIG
				case "HKEY_PERFORMANCE_DATA":
					WorkingRegPath.RootKey = registry.PERFORMANCE_DATA
				default:
					log.Printf("unknown: %#v\n", rootkeySplit)
					line = ""
					continue
				}

				if WorkingRegPath.OpenHandle {
					line = ""
					WorkingRegPath.Close()
				}

				if err := WorkingRegPath.Open(); err != nil {
					if err == registry.ErrNotExist {
						WorkingRegPath.InfoHeader = true
						WorkingRegPath.Log(Err, "\n["+line+"]")
					} else {
						WorkingRegPath.InfoHeader = true
						WorkingRegPath.Log(Panic, "\n["+line+"] ("+err.Error()+")")
					}

					line = ""
					continue
				}

				line = ""

			} else {
				if !WorkingRegPath.OpenHandle {
					if WorkingRegPath.OpenHandle {
						WorkingRegPath.Log(Warn, line)
					} else {
						WorkingRegPath.Log(Err, line)
					}

					line = ""
					continue
				}

				equalPos := strings.Index(line, "=")

				fileData := FileData{
					Key:   strings.Trim(line[:equalPos], `"`),
					Value: line[equalPos+1:],
				}
				fileData.Type = GetClassification(fileData.Value)

				switch fileData.Type {
				case registry.SZ:
					fileData.Value = strings.Trim(fileData.Value, `"`)
					fileData.Value = strings.ReplaceAll(fileData.Value, `\\`, `\`) // TODO: escape func?

					regValue, _, err := WorkingRegPath.RegPath.GetStringValue(fileData.Key)
					if err == registry.ErrNotExist {
						WorkingRegPath.Log(Err, line) // Entry not found

						// fileData.Clear()
						line = ""
						continue
					} else if err == registry.ErrUnexpectedType {
						WorkingRegPath.Log(Err, line+" (unexpected key value type)") // Entry not in the same format

						// fileData.Clear()
						line = ""
						continue
					} else if err != nil { // Error
						log.Println(err)
						log.Fatal(line)
					}

					if regValue == fileData.Value {
						WorkingRegPath.Log(Info, line)
					} else {
						WorkingRegPath.Log(Warn, line)
					}

				case registry.DWORD:
					RegFileValue := strings.TrimLeft(fileData.Value[6:], "0") // Deletes the dword: and all preceding 0s

					regValue, _, err := WorkingRegPath.RegPath.GetIntegerValue(fileData.Key)
					if err == registry.ErrNotExist {
						WorkingRegPath.Log(Err, line) // Entry not found

						// fileData.Clear()
						line = ""
						continue
					} else if err == registry.ErrUnexpectedType {
						WorkingRegPath.Log(Err, line+" (unexpected key value type)") // Entry not in the same format

						// fileData.Clear()
						line = ""
						continue
					} else if err != nil { // Error
						log.Println(err)
						log.Fatal(line)
					}

					if CompareDWord(DecToHexstring(regValue), RegFileValue) {
						WorkingRegPath.Log(Info, line)
					} else {
						WorkingRegPath.Log(Warn, line)
					}

				case registry.EXPAND_SZ: // REG_EXPAND_SZ

					regValue, _, err := WorkingRegPath.RegPath.GetStringValue(fileData.Key)
					if err == registry.ErrNotExist {
						WorkingRegPath.Log(Err, line) // Entry not found

						// fileData.Clear()
						line = ""
						continue
					} else if err == registry.ErrUnexpectedType {
						WorkingRegPath.Log(Err, line+" (unexpected key value type)") // Entry not in the same format

						// fileData.Clear()
						line = ""
						continue
					} else if err != nil { // Error
						log.Println(err)
						log.Fatal(line)
					}

					regValue = strings.ReplaceAll(regValue, `\\`, `\`) // escape func?

					if regValue == BinstringToString(fileData.Value[7:]) {
						// WorkingRegPath.Log(Info, fileData.Key+"="+fileData.Value)
						WorkingRegPath.Log(Info, line)
					} else {
						// WorkingRegPath.Log(Warn, fileData.Key+"="+fileData.Value)
						WorkingRegPath.Log(Warn, line)
					}

				case registry.BINARY: // Binär
					log.Println("// TODO: " + fileData.Key + "=" + fileData.Value)
					line = ""
				default:
					log.Println("// TODO: unknown typ", fileData)
				}

				// fileData.Clear()
				line = ""
			}

		} else {
			line = ""
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s\n", err)
	}

	fmt.Println("Press 'Enter' to exit...")
	fmt.Scanln()
}

func trimSpace(r rune) bool {
	// return !unicode.IsLetter(r) && !unicode.IsNumber(r) // TODO: Testing
	return unicode.IsSpace(r)
}
