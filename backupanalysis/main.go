package main

import (
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/spf13/pflag"
	"github.com/willauld/utils/hashdictionary"
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
	Name     string
	Path     string
	Md5      string
	Size     int64
	Modified time.Time
}

type dataRecord struct {
	SourcePath string
	TargetPath string
	MasterList map[string]fileInfo
}

var dr = dataRecord{MasterList: map[string]fileInfo{}}

//var dr dataRecord

type payload struct {
	fs []*fileInfo
}

func value2payload(v *hashdictionary.ValuePayload) (i *payload) {
	return (*payload)(unsafe.Pointer(v))
}
func payload2value(p *payload) (v *hashdictionary.ValuePayload) {
	return (*hashdictionary.ValuePayload)(unsafe.Pointer(p))
}

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

		dr.MasterList[path] =
			fileInfo{f.Name(), path, md5str, f.Size(), f.ModTime()}
		//fmt.Printf("\t%s %s ::%v\n", md5str, f.Name(), f.ModTime())
	}
	return nil
}

var dictCount int

func display(a hashdictionary.DictEntry) {
	fmt.Printf("%d: %s:\n", dictCount, a.Str)
	if a.Value != nil {
		for i, fi := range value2payload(a.Value).fs {
			fmt.Printf("\t%d - %d bytes %s\n", i, fi.Size, fi.Path)
		}
		dictCount++
	}
}

func masterToMD5Dict() {
	for _, v := range dr.MasterList {
		te, err := md5Dict.Insert(hashdictionary.DictEntry{Str: v.Md5})
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
		newCopy := v
		pl.fs = append(pl.fs, &newCopy)
	}
}

func storeToFile(data *dataRecord, dataLog string) error {
	// serialize the data
	// open and write to file ##### need to update with multiple source sets
	dataFile, err := os.Create(dataLog)
	if err != nil {
		fmt.Printf("could not open DataLog: %s :: %s\n", dataLog, err)
		return err
	}
	defer dataFile.Close()
	encoder := gob.NewEncoder(dataFile)
	err = encoder.Encode(data)
	return err
}

func loadFromFile(dataLog string) (data dataRecord, err error) {
	// open the data file
	dataFile, err := os.Open(dataLog)
	if err != nil {
		log.Fatalf("DataLog [%s] Open Error: %v\n", dataLog, err)
		return data, err
	}
	defer dataFile.Close()

	decoder := gob.NewDecoder(dataFile)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Printf("decode() Error: %v\n", err)
	}
	return data, err
}

func main() {
	versionPtr := pflag.Bool("version", false, "program version")
	srcPtr := pflag.String("src", "", "source directory")
	dataPtr := pflag.String("data", "", "data file")
	pflag.Parse()
	//fmt.Println("input:", *srcPtr)
	//fmt.Println("data:", *dataPtr)
	//fmt.Println("tail:", pflag.Args())

	if *versionPtr == true {
		fmt.Printf("\t Version %d.%d", version.major, version.minor)
		os.Exit(0)
	}
	dataLog := *dataPtr
	if dataLog == "" {
		//dataLog = "~/.backupAnalysis.dat"
		dataLog = "c:/home/auld/.backupAnalysis.dat"
		//TODO need to make this look up the users home dir and put this file there
	}
	sourceDir := *srcPtr
	if sourceDir == "" {
		fmt.Printf("sourceDir is the empty string\n")
		sourceDir = "c:/home/auld/temp/backupCopyTest"
	}

	fmt.Println("Source: " + sourceDir)
	fmt.Println("data: " + dataLog)

	// check if the source dir exist
	src, err := os.Stat(sourceDir)
	if err != nil {
		panic(err)
	}

	if !src.IsDir() {
		fmt.Println("Source is not a directory")
		os.Exit(1)
	}

	ldr, err := loadFromFile(dataLog)
	if err == nil {
		// OK if no file found to load from
		//fmt.Printf("####\n%+v\n#####\n", ldr)
		dr = ldr
	}

	err = filepath.Walk(sourceDir, a)
	if err != nil {
		fmt.Printf("filepath.Walk failed %v\n", err)
		return
	}
	fmt.Printf("%d files processed\n", len(dr.MasterList))

	err = storeToFile(&dr, dataLog)
	if err != nil {
		fmt.Printf("storeToFile failed: %v\n", err)
	}

	//fmt.Printf("####\n%+v\n#####\n", dr)

	md5Dict = hashdictionary.Create(hashSize, hash, equal)
	masterToMD5Dict()
	md5Dict.Map(display)
}
