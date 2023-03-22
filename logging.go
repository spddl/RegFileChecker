package main

import (
	"fmt"
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
		fmt.Printf(WithoutColor, "\n["+FindRootKey[d.RootKey]+"\\"+d.StringPath+"]")
	}
}

func (d *Data) Log(level LogLevel, line string) {
	switch level {
	case Info:
		if !NoInfo {
			d.LogInfoHeader()
			if NoColor {
				fmt.Printf(WithoutColor, line)
			} else {
				fmt.Printf(InfoColor, line)
			}
		}
	case Warn:
		if !NoWarn {
			d.LogInfoHeader()
			if NoColor {
				fmt.Printf(WithoutColor, line)
			} else {
				fmt.Printf(WarningColor, line)
			}
		}
		issues++
	case Err:
		if !NoErr {
			d.LogInfoHeader()
			if NoColor {
				fmt.Printf(WithoutColor, line)
			} else {
				fmt.Printf(ErrorColor, line)
			}
		}
		issues++
	case Panic:
		d.LogInfoHeader()
		if NoColor {
			fmt.Printf(WithoutColor, line)
		} else {
			fmt.Printf(PanicColor, line)
		}
		issues++
	}
}
