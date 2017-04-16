package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/willauld/utils/hashdictionary"
)

func init() {
	createSource()
}

func createSource() {

	fmt.Println("Golang Version:", runtime.Version())

	src := os.TempDir()
	fmt.Printf("TempDir: %s\n", src)
}

func deleteSource() {

}

func deleteTarget() {

}

func TestHash(t *testing.T) {
	tstCases := []struct {
		in  hashdictionary.DictEntry
		out int
	}{
		{hashdictionary.DictEntry{Str: ""}, 0},
		{hashdictionary.DictEntry{Str: "testItem1"}, 375},
		{hashdictionary.DictEntry{Str: "testItem2"}, 376},
		{hashdictionary.DictEntry{Str: "a"}, 97},
	}
	for _, tstCase := range tstCases {
		rtnv := hash(tstCase.in)
		if rtnv != tstCase.out {
			t.Errorf("hash(%s) expected %d but returned %d\n",
				tstCase.in.Str, tstCase.out, rtnv)
		}
	}
}

func TestEqual(t *testing.T) {
	tstCases := []struct {
		in1 hashdictionary.DictEntry
		in2 hashdictionary.DictEntry
		out bool
	}{
		{hashdictionary.DictEntry{Str: ""}, hashdictionary.DictEntry{Str: ""}, true},
		{hashdictionary.DictEntry{Str: "a"}, hashdictionary.DictEntry{Str: ""}, false},
		{hashdictionary.DictEntry{Str: "a"}, hashdictionary.DictEntry{Str: "b"}, false},
		{hashdictionary.DictEntry{Str: "a"}, hashdictionary.DictEntry{Str: "b"}, false},
		{hashdictionary.DictEntry{Str: "this is a bigger string"}, hashdictionary.DictEntry{Str: "this is a bigger string"}, true},
	}
	for _, tstCase := range tstCases {
		rtnv := equal(tstCase.in1, tstCase.in2)
		if rtnv != tstCase.out {
			t.Errorf("equal(%s, %s) expected %t but returned %t\n",
				tstCase.in1.Str, tstCase.in2.Str, tstCase.out, rtnv)
		}
	}
}

func TestA(t *testing.T) {
	fmt.Printf("TestA WIP: work in progress\n")
}

func TestDisplay(t *testing.T) {
	fmt.Printf("TestDisplay WIP: work in progress\n")

}

func buildDataRecord(s string, llen int) *dataRecord {
	mydr := dataRecord{s,
		"no target", map[string]fileInfo{},
	}
	fmod := time.Date(1960, 12, 1, 0, 0, 0, 0, time.Local)
	for i := 0; i < llen; i++ {
		name := fmt.Sprintf("name%d", i)
		path := fmt.Sprintf("/path/name%d", i)
		md5 := fmt.Sprintf("md5_%d", i)
		size := int64(30 + i)
		mydr.MasterList[path] =
			fileInfo{name, path, md5, size, fmod}
		//mydr.MasterList = append(mydr.MasterList,
		//	fileInfo{name, path, md5, size, fmod})
	}
	return &mydr
}

func TestLoadStoreGob(t *testing.T) {
	sourceDir := "c:/meme"
	numFileInfo := 2
	mydr := buildDataRecord(sourceDir, numFileInfo)
	fmt.Printf("%+v\n", mydr)

	src := os.TempDir()
	dataLog := src + "/.backupAnalysis.dat"
	err := storeToFile(mydr, dataLog)
	if err != nil {
		t.Errorf("storeToFile() failed: %s\n", err)
	} else {
		dr, err := loadFromFile(dataLog)
		if err != nil {
			t.Errorf("loadFromFile() failed: %s\n", err)
		} else {
			fmt.Printf("FromFile:\n%+v\n", mydr)
			if dr.SourcePath != sourceDir {
				t.Errorf("dr.source expected to be [%s] but is [%s]\n",
					sourceDir, dr.SourcePath)
			}
			if dr.TargetPath != "no target" {
				t.Errorf("dr.TargetPath expected to be [%s] but is [%s]\n",
					"no target", dr.TargetPath)
			}
			if len(dr.MasterList) != numFileInfo {
				t.Errorf("dr.masterList should have %d entry but has %d\n",
					numFileInfo, len(dr.MasterList))
			}
		}
	}

	err = os.Remove(dataLog)
	if err != nil {
		t.Errorf("Remove() failed: %s\n", err)
	}
}

// TEST GOB standalone
type j struct {
	A, b int
}
type User struct {
	Name, Pass string //fields mast be exportable to be transfered by gob
	List       []j
}

func TestGob(t *testing.T) {
	src := os.TempDir()
	file := src + "/test.gob"

	var datato = &User{Name: "Donald", Pass: "DuckPass", List: []j{{1, 2}, {2, 1}}}
	var datafrom = new(User)

	err := Save(file, datato)
	Check(err)
	err = Load(file, datafrom)
	Check(err)
	fmt.Println(datafrom)

	err = os.Remove(file)
	Check(err)
}

// Encode via Gob to file
func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		defer file.Close()
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(object)
	}
	return err
}

// Decode Gob file
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	return err
}

func Check(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(line, "\t", file, "\n", e)
		os.Exit(1)
	}
}
