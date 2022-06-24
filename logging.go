package main

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

type LogLevel int

const (
	Info LogLevel = iota
	Warn
	Err
	Panic
)

func (d *Data) LogInfoHeader() {
	if !d.InfoHeader {
		d.InfoHeader = true
		fmt.Printf(NoColor, "\n["+FindRootKey(d.RootKey)+"\\"+d.StringPath+"]")
	}
}

func (d *Data) Log(level LogLevel, line string) {
	switch level {
	case Info:
		if !NoInfo {
			d.LogInfoHeader()
			fmt.Printf(InfoColor, line)
		}
	case Warn:
		if !NoWarn {
			d.LogInfoHeader()
			fmt.Printf(WarningColor, line)
		}
	case Err:
		if !NoErr {
			d.LogInfoHeader()
			fmt.Printf(ErrorColor, line)
		}
	case Panic:
		d.LogInfoHeader()
		fmt.Printf(PanicColor, line)
	}
}

func FindRootKey(rootKey registry.Key) string {
	switch rootKey {
	case registry.CLASSES_ROOT:
		return "HKEY_CLASSES_ROOT"
	case registry.CURRENT_USER:
		return "HKEY_CURRENT_USER"
	case registry.LOCAL_MACHINE:
		return "HKEY_LOCAL_MACHINE"
	case registry.USERS:
		return "HKEY_USERS"
	case registry.CURRENT_CONFIG:
		return "HKEY_CURRENT_CONFIG"
	case registry.PERFORMANCE_DATA:
		return "HKEY_PERFORMANCE_DATA"
	default:
		return ""
	}
}
