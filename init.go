package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"unicode"

	"golang.org/x/sys/windows"
)

var (
	NoInfo       bool
	NoWarn       bool
	NoErr        bool
	NoColor      bool
	Exit         bool
	Verbose      bool
	RegFilePaths []string
)

const (
	WithoutColor = "%s\n"
	InfoColor    = "\033[1;32m%s\033[0m\n"
	WarningColor = "\033[1;33m%s\033[0m\n"
	ErrorColor   = "\033[1;31m%s\033[0m\n"
	PanicColor   = "\033[1;35m%s\033[0m\n"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.BoolVar(&NoInfo, "noinfo", false, "Without info")
	flag.BoolVar(&NoWarn, "nowarn", false, "Without warnung")
	flag.BoolVar(&NoErr, "noerr", false, "Without error")
	flag.BoolVar(&NoColor, "nocolor", false, "Without color")
	flag.BoolVar(&Exit, "exit", false, "Exit after processing")

	flag.BoolVar(&Verbose, "verbose", false, "Verbose line")
	flag.Parse()

	separatorFunc := func(c rune) bool {
		return unicode.IsSpace(c) || c == ';' || c == ','
	}

	args := flag.Args()
	if len(args) != 0 {
		RegFilePaths = strings.FieldsFunc(args[0], separatorFunc)
	}

	enableColor()
}

func enableColor() {
	var handle = windows.Handle(os.Stdout.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(handle, &mode); err == nil {
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING

		windows.SetConsoleMode(handle, mode)
	}
}
