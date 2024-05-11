package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding/unicode"
)

// http://www.sentinelchicken.com/data/TheWindowsNTRegistryFileFormat.pdf

type FileData struct {
	Value string
	Key   string
	Type  uint32
}

var issues int

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	for _, RegFilePath := range RegFilePaths {
		file, err := os.Open(RegFilePath)
		if err != nil {
			log.Println(err)
			continue
		}

		var scanner *bufio.Scanner
		contentType := DetectContentType(file)
		file.Seek(0, 0)
		switch contentType {
		case "utf-8":
			scanner = bufio.NewScanner(file)
		case "utf-16le":
			scanner = bufio.NewScanner(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder().Reader(file))
		default:
			file.Close()
			fmt.Fprintln(os.Stderr, "ContentType:", contentType)
			os.Exit(1)
		}

		var WorkingRegPath = new(Data)
		var index uint
		var line string
		for scanner.Scan() {
			index++
			if index == 1 {
				continue // skip "Windows Registry Editor Version"
			}

			if rawLine := scanner.Text(); len(rawLine) > 1 {
				// fmt.Printf("[%d] %s\n", index, rawLine)
				if Verbose {
					fmt.Printf("[%d] %s\n", index, rawLine)
				}

				commentIndex := strings.Index(rawLine, ";")
				if commentIndex != -1 {
					rawLine = rawLine[:commentIndex] // remove comments
				}

				if strings.HasSuffix(rawLine, "\\") {
					line += strings.TrimSuffix(strings.TrimSpace(rawLine), "\\")
					continue
				}

				line += strings.TrimSpace(rawLine) // remove spaces
				if line == "" {
					continue // skip empty lines
				}

				if line[0] == 91 { // [
					WorkingRegPath = new(Data)
					line = line[1 : len(line)-1]

					var rootkeySplit = strings.SplitN(line, `\`, 2)
					WorkingRegPath.StringPath = rootkeySplit[1]

					if rootkeySplit[0][0] == 45 { // [-HKEY_LOCAL_MACHINE\...
						WorkingRegPath.DelRootKey = true
						rootkeySplit[0] = rootkeySplit[0][1:] // remove "-"
					}

					if val, ok := FindRootString[rootkeySplit[0]]; ok {
						WorkingRegPath.RootKey = val
					} else {
						line = ""
						log.Panicf("unknown: %#v\n", rootkeySplit)
					}

					if WorkingRegPath.OpenHandle {
						line = ""
						WorkingRegPath.Close()
					}

					if err := WorkingRegPath.Open(); err != nil {
						switch {
						case err == registry.ErrNotExist && WorkingRegPath.DelRootKey:
							WorkingRegPath.InfoHeader = true
							WorkingRegPath.Log(Info, "\n["+line+"]")
						case err == registry.ErrNotExist:
							WorkingRegPath.InfoHeader = true
							WorkingRegPath.Log(Err, "\n["+line+"] ;"+err.Error())
						default:
							WorkingRegPath.InfoHeader = true
							WorkingRegPath.Log(Panic, "\n["+line+"] ;"+err.Error())
						}
						line = ""
						continue
					}

					if WorkingRegPath.DelRootKey { // Path exists
						WorkingRegPath.InfoHeader = true
						WorkingRegPath.Log(Err, "\n["+line+"]")
						line = ""
						WorkingRegPath.Close()
						continue
					}

					line = ""
				} else {
					if !WorkingRegPath.OpenHandle {
						WorkingRegPath.Log(Err, line)

						line = ""
						continue
					}

					equalPos := strings.Index(line, "=")
					if equalPos == -1 {
						line = ""
						continue
					}
					fileData := FileData{
						Key:   strings.Trim(line[:equalPos], `"`),
						Value: line[equalPos+1:],
					}
					fileData.Type = GetClassification(fileData.Value)
					if fileData.Key == "@" {
						fileData.Key = "" // default value
					}

					switch fileData.Type {
					case registry.NONE:
						_, _, err := WorkingRegPath.RegPath.GetValue(fileData.Key, nil)
						if err == registry.ErrNotExist {
							WorkingRegPath.Log(Info, line)
						} else {
							WorkingRegPath.Log(Err, line)
						}
					case registry.SZ:
						fileData.Value = strings.ReplaceAll(fileData.Value, `\"`, `"`)
						fileData.Value = strings.ReplaceAll(fileData.Value, `\\`, `\`)
						if fileData.Value[0] == 34 && fileData.Value[len(fileData.Value)-1] == 34 {
							fileData.Value = fileData.Value[1 : len(fileData.Value)-1]
						}

						regValue, _, err := WorkingRegPath.RegPath.GetStringValue(fileData.Key)
						switch err {
						case nil: // none
						case registry.ErrNotExist:
							WorkingRegPath.Log(Err, line) // Entry not found
							line = ""
							continue
						case registry.ErrUnexpectedType:
							WorkingRegPath.Log(Err, line+" ;unexpected key value type") // Entry not in the same format
							line = ""
							continue
						default: // Error
							log.Println(err)
							log.Fatal(line)
						}

						if regValue == fileData.Value {
							WorkingRegPath.Log(Info, line)
						} else {
							WorkingRegPath.Log(Warn, line+" ;current: "+regValue)
						}

					case registry.DWORD:
						RegFileValue := strings.TrimLeft(fileData.Value[6:], "0") // Deletes the dword: and all preceding 0s

						regValue, _, err := WorkingRegPath.RegPath.GetIntegerValue(fileData.Key)
						switch err {
						case nil: // none
						case registry.ErrNotExist:
							WorkingRegPath.Log(Err, line) // Entry not found
							line = ""
							continue
						case registry.ErrUnexpectedType:
							WorkingRegPath.Log(Err, line+" ;unexpected key value type") // Entry not in the same format
							line = ""
							continue
						default:
							log.Println(err)
							log.Fatal(line)
						}

						if CompareDWord(DecToHexstring(regValue), RegFileValue) {
							WorkingRegPath.Log(Info, line)
						} else {
							WorkingRegPath.Log(Warn, line+" ;current: 0x"+DecToHexstring(regValue))
						}

					case registry.EXPAND_SZ: // REG_EXPAND_SZ
						regValue, _, err := WorkingRegPath.RegPath.GetStringValue(fileData.Key)
						switch err {
						case nil: // none
						case registry.ErrNotExist:
							WorkingRegPath.Log(Err, line) // Entry not found
							line = ""
							continue
						case registry.ErrUnexpectedType:
							WorkingRegPath.Log(Err, line+" ;unexpected key value type") // Entry not in the same format
							line = ""
							continue
						default:
							log.Println(err)
							log.Fatal(line)
						}

						regValue = strings.ReplaceAll(regValue, `\\`, `\`)     // escape func?
						if regValue == BinstringToString(fileData.Value[7:]) { // without "hex(2):"
							WorkingRegPath.Log(Info, line)
						} else {
							WorkingRegPath.Log(Warn, line+" ;current: "+regValue)
						}

					case registry.BINARY: // Binär
						regValue, _, err := WorkingRegPath.RegPath.GetBinaryValue(fileData.Key)
						switch err {
						case nil: // none
						case registry.ErrNotExist:
							WorkingRegPath.Log(Err, line) // Entry not found
							line = ""
							continue
						case registry.ErrUnexpectedType:
							WorkingRegPath.Log(Err, line+" ;unexpected key value type") // Entry not in the same format
							line = ""
							continue
						default:
							log.Println(err)
							log.Fatal(line)
						}

						data, err := hex.DecodeString(strings.ReplaceAll(fileData.Value[strings.Index(fileData.Value, ":")+1:], ",", ""))
						if err != nil {
							WorkingRegPath.Log(Err, line) // Entry not in the same format
							line = ""
							continue
						}

						if bytes.EqualFold(regValue, data) {
							WorkingRegPath.Log(Info, line)
						} else {
							WorkingRegPath.Log(Warn, line)
						}

						line = ""

					case registry.QWORD: // QWORD
						fallthrough
					case registry.MULTI_SZ: // MULTI_SZ
						regValue, _, err := WorkingRegPath.RegPath.GetIntegerValue(fileData.Key)
						switch err {
						case nil: // none
						case registry.ErrNotExist:
							WorkingRegPath.Log(Err, line) // Entry not found
							line = ""
							continue
						case registry.ErrUnexpectedType:
							WorkingRegPath.Log(Err, line+" ;unexpected key value type") // Entry not in the same format

							line = ""
							continue
						default:
							log.Println(err)
							log.Fatal(line)
						}

						b := BinStringToByteArray(fileData.Value[7:])
						i := ByteArrayToUint64(b)

						if i == regValue {
							WorkingRegPath.Log(Info, line)
						} else {
							WorkingRegPath.Log(Warn, fmt.Sprintf("%s ;current: %s", line, ByteArrayToBinString(Uint64ToByteArray(regValue))))
						}

						line = ""
					default:
						log.Printf("// TODO: unknown typ: %#v\n", fileData)
						continue
					}

					line = ""
				}
			} else {
				line = ""
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Invalid input: %s\n", err)
		}
		file.Close()
	}

	if issues != 0 {
		fmt.Printf("; issues: %d\n; %s\n", issues, strings.Join(RegFilePaths, ","))
	}

	if !Exit {
		fmt.Printf("\nPress 'Enter' to exit...\n")
		fmt.Scanln()
	}
}

func DetectContentType(file *os.File) string {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		log.Panicln("Error:", err)
	}
	contentType := http.DetectContentType(buffer[:n])
	index := strings.Index(contentType, "; charset=")
	if index == -1 {
		return ""
	}
	return contentType[index+10:]
}
