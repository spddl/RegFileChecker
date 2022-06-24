package main

import (
	"flag"
)

var (
	NoInfo      bool
	NoWarn      bool
	NoErr       bool
	RegFilePath string
)

const (
	NoColor      = "%s\n"
	InfoColor    = "\033[1;32m%s\033[0m\n"
	WarningColor = "\033[1;33m%s\033[0m\n"
	ErrorColor   = "\033[1;31m%s\033[0m\n"
	PanicColor   = "\033[1;35m%s\033[0m\n"
)

func init() {
	flag.BoolVar(&NoInfo, "noinfo", false, "Without info")
	flag.BoolVar(&NoWarn, "nowarn", false, "Without warnung")
	flag.BoolVar(&NoErr, "noerr", false, "Without error")
	flag.Parse()
	RegFilePath = flag.Args()[0]
}
