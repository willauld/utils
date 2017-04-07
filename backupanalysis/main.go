package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/spf13/pflag"
	"github.com/willauld/backuputils/hashdictionary"
)

var version = struct {
	major int
	minor int
}{0, 1}
var (
	force     bool
	newTarget bool
)

var md5Dict *hashdictionary.HashDictionary

// hashSize should be prime. Some choices are:
//   163, 317, 631, 1259, 2503, 9973 # 23 for testing
const hashSize = 631

// IDEA HERE is to have the file info I want to track cept in the
//fileInfo slice and pass the payload field to hashDictionare to
// track it with the md5str. For this to work I will NEED to be
// garenteed that the offset from payload to payload.i is an unchanging
// constant. This can be check with the unsafe package.
type fileInfo struct {
	name string
	path string
	md5  string
	size int64
}

type payload struct {
	fs []fileInfo
}

func value2payload(v *hashdictionary.ValuePayload) (i *payload) {
	return (*payload)(unsafe.Pointer(v))
}
func payload2value(p *payload) (v *hashdictionary.ValuePayload) {
	return (*hashdictionary.ValuePayload)(unsafe.Pointer(p))
}

/*
func outer2inner(o *payload) (i hashdictionary.DictEntry) {
	return *(*hashdictionary.DictEntry)(unsafe.Pointer(&o.i))
}

func inner2outer(i hashdictionary.DictEntry) (o *payload) {
	s := payload{}
	return unsafe.Pointer(uintptr(unsafe.Pointer(&i)) - unsafe.Offsetof(s.i))
}*/

func hash(item hashdictionary.DictEntry) int {
	str := item.Str
	v := 0
	for i := 0; i < len(str); i++ {
		v = (10*v + int(str[i])) % hashSize
	}
	return v
}

func equal(item1, item2 hashdictionary.DictEntry) bool {
	str1 := item1.Str
	str2 := item2.Str
	return str1 == str2
}

func a(path string, f os.FileInfo, err error) error {
	//fmt.Printf("Walking: %s\n", path)
	if !f.IsDir() {
		fil, ferr := os.Open(path)
		if ferr != nil {
			log.Fatal(ferr)
		}
		defer fil.Close()

		h := md5.New()
		if _, ferr := io.Copy(h, fil); ferr != nil {
			log.Fatal(ferr)
		}
		md5str := fmt.Sprintf("%x", h.Sum(nil))
		te, err := md5Dict.Insert(hashdictionary.DictEntry{Str: md5str})
		if err != nil {
			log.Fatal(err)
		}
		var pl *payload
		if te.Value == nil {
			pl = &payload{}
			te.Value = payload2value(pl)
		} else {
			pl = value2payload(te.Value)
		}
		pl.fs = append(pl.fs, fileInfo{f.Name(), path, md5str, f.Size()})

		//fmt.Printf("\t%s %s\n", md5str, f.Name())
	}
	return nil
}

var dictCount int

func display(a hashdictionary.DictEntry) {
	fmt.Printf("%d: %s:\n", dictCount, a.Str)
	if a.Value != nil {
		for i, fi := range value2payload(a.Value).fs {
			fmt.Printf("\t%d - %d bytes %s\n", i, fi.size, fi.path)
		}
		dictCount++
	}
}

func main() {
	versionPtr := pflag.Bool("version", false, "program version")
	srcPtr := pflag.String("src", "", "source directory")
	pflag.Parse()
	fmt.Println("input:", *srcPtr)
	fmt.Println("tail:", pflag.Args())

	if *versionPtr == true {
		fmt.Printf("\t Version %d.%d", version.major, version.minor)
		os.Exit(0)
	}
	sourceDir := *srcPtr

	if sourceDir == "" {
		fmt.Printf("sourceDir is the empty string\n")
		sourceDir = "c:/home/auld/temp/backupCopyTest"
	}

	fmt.Println("Source :" + sourceDir)

	// check if the source dir exist
	src, err := os.Stat(sourceDir)
	if err != nil {
		panic(err)
	}

	if !src.IsDir() {
		fmt.Println("Source is not a directory")
		os.Exit(1)
	}

	md5Dict = hashdictionary.Create(hashSize, hash, equal)

	err = filepath.Walk(sourceDir, a)
	if err != nil {
		fmt.Printf("filepath.Walk failed %v\n", err)
		return
	}
	md5Dict.Map(display)
}
