package main

import (
	"log"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type Data struct {
	InfoHeader bool
	OpenHandle bool
	RootKey    registry.Key
	RegPath    registry.Key
	StringPath string
}

func (d *Data) Open() error {
	var err error
	d.RegPath, err = registry.OpenKey(d.RootKey, d.StringPath, registry.QUERY_VALUE)
	if err != nil {
		return err
	}

	d.OpenHandle = true
	return nil
}

func (d *Data) Close() {
	if !d.check() {
		return
	}

	err := d.RegPath.Close()
	if err != nil {
		log.Printf("%#v\n", d)
		// log.Fatal(err)
	}
	d.OpenHandle = false
}

func GetClassification(value string) int {
	switch {
	case strings.HasPrefix(value, `"`) || strings.HasPrefix(value, "@"):
		return registry.SZ
	case strings.HasPrefix(value, "hex:"):
		return registry.BINARY
	case strings.HasPrefix(value, "dword:"):
		return registry.DWORD
	case strings.HasPrefix(value, "hex(0):"):
		return registry.NONE
	case strings.HasPrefix(value, "hex(1):"):
		return registry.SZ
	case strings.HasPrefix(value, "hex(2):"):
		return registry.EXPAND_SZ
	case strings.HasPrefix(value, "hex(3):"):
		return registry.BINARY
	case strings.HasPrefix(value, "hex(4):"):
		return registry.DWORD
	case strings.HasPrefix(value, "hex(5):"):
		return registry.DWORD_BIG_ENDIAN
	case strings.HasPrefix(value, "hex(7):"):
		return registry.MULTI_SZ
	case strings.HasPrefix(value, "hex(8):"):
		return registry.RESOURCE_LIST
	case strings.HasPrefix(value, "hex(a):"):
		return registry.RESOURCE_REQUIREMENTS_LIST
	case strings.HasPrefix(value, "hex(b):"):
		return registry.QWORD
	}
	return 0
}

func (d *Data) check() bool {
	if d.RootKey == registry.Key(0) ||
		d.RegPath == registry.Key(0) ||
		d.StringPath == "" {
		return false
	}
	return true
}
