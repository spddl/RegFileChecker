package main

import (
	"log"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type Data struct {
	InfoHeader bool
	OpenHandle bool
	DelRootKey bool
	RootKey    registry.Key
	RegPath    registry.Key
	StringPath string
}

var FindRootKey = map[registry.Key]string{
	registry.CLASSES_ROOT:     "HKEY_CLASSES_ROOT",
	registry.CURRENT_USER:     "HKEY_CURRENT_USER",
	registry.LOCAL_MACHINE:    "HKEY_LOCAL_MACHINE",
	registry.USERS:            "HKEY_USERS",
	registry.CURRENT_CONFIG:   "HKEY_CURRENT_CONFIG",
	registry.PERFORMANCE_DATA: "HKEY_PERFORMANCE_DATA",
}

var FindRootString = map[string]registry.Key{
	"HKEY_CLASSES_ROOT":     registry.CLASSES_ROOT,
	"HKEY_CURRENT_USER":     registry.CURRENT_USER,
	"HKEY_LOCAL_MACHINE":    registry.LOCAL_MACHINE,
	"HKEY_USERS":            registry.USERS,
	"HKEY_CURRENT_CONFIG":   registry.CURRENT_CONFIG,
	"HKEY_PERFORMANCE_DATA": registry.PERFORMANCE_DATA,
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
	}
	d.OpenHandle = false
}

func GetClassification(value string) uint32 {
	switch {
	case value == "-":
		return registry.NONE

	case strings.HasPrefix(value, `"`):
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
