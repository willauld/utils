package hashdictionary

import (
	"fmt"
	"testing"
)

//type DictEntry string

// hashSize should be prime. Some choices are:
//   163, 317, 631, 1259, 2503, 9973 # 23 for testing
const hashSize = 23

func hash(item DictEntry) int {
	str := item.Str
	v := 0
	for i := 0; i < len(str); i++ {
		v = (10*v + int(str[i])) % hashSize
	}
	return v
}

func equal(item1, item2 DictEntry) bool {
	str1 := item1.Str
	str2 := item2.Str
	return str1 == str2
}

func TestBasic(t *testing.T) {
	tstCases := []struct {
		line     DictEntry
		expected DictEntry
	}{
		{DictEntry{Str: "testItem1"}, DictEntry{Str: "testItem1"}},
		{DictEntry{Str: "testItem2"}, DictEntry{Str: "testItem2"}},
	}
	mydict := Create(hashSize, hash, equal)

	for _, tstCase := range tstCases {
		_, err := mydict.Insert(tstCase.line)
		if err != nil {
			t.Errorf("70:err is non-nil on Insert operation, always unexpected: %v\n", err)
			//panic(err)
		}
		resultEntry, err := mydict.Lookup(tstCase.line)
		if err != nil {
			t.Errorf("60: err is non-nil on lookup operation: %v\n", err)
			//panic(err)
		}
		if resultEntry == nil {
			t.Errorf("20:unexpected resultEntry returned as nil for lookup of %s\n", tstCase.line.Str)
		}
		if resultEntry.Str != tstCase.line.Str {
			t.Errorf("50:lookup(%s) returned [%s] expected [%s]\n", tstCase.line.Str, resultEntry.Str, tstCase.line.Str)
		}
		err = mydict.Delete(tstCase.line)
		if err != nil {
			t.Errorf("40:err is non-nil on delete operation: %v\n", err)
			//panic(err)
		}
		resultEntry, err = mydict.Lookup(tstCase.line)
		if err == nil {
			t.Errorf("30:err is nil on lookup operation that should not return nil\n")
			//panic(err)
		}
		if resultEntry != nil {
			t.Errorf("10: non-nil resultEntry returned as lookup of non-existant entry: %s\n", tstCase.line.Str)
		}
		err = mydict.Delete(tstCase.line)
		if err == nil {
			t.Errorf("err is nil on delete operation that should not return nil\n")
			panic(err) // Delete should fail when call a second time.
		}
	}
}

func loadTable() *HashDictionary {

	entries := []string{
		"a", "b", "c", "d",
		"e", "f", "g", "h",
		"i", "j", "k", "l",
		"m", "n", "o", "p",
		"q", "r", "s", "t",
		"u", "v", "w", "x",
		"y", "z", "ww", "xx",
		"A", "B", "C", "D",
		"E", "F", "G", "H",
		"I", "J", "K", "L",
		"M", "N", "O", "P",
		"Q", "R", "S", "T",
		"U", "V", "W", "X",
		"Y", "Z", "Ww", "Xx",
		"ja", "jb", "jc", "jd",
		"je", "jf", "jg", "jh",
		"ji", "jj", "jk", "jl",
		"jm", "jn", "jo", "jp",
		"jq", "jr", "js", "jt",
		"ju", "jv", "jw", "jx",
		"jy", "jz", "jww", "jxx",
		"jA", "jB", "jC", "jD",
		"jE", "jF", "jG", "jH",
		"jI", "jJ", "jK", "jL",
		"jM", "jN", "jO", "jP",
		"jQ", "jR", "jS", "jT",
		"jU", "jV", "jW", "jX",
		"jY", "jZ", "jWw", "jXx",
	}

	mydict := Create(hashSize, hash, equal)
	// set up initial dictionary entries
	for _, tstCase := range entries {
		_, err := mydict.Insert(DictEntry{Str: tstCase})
		if err != nil {
			panic(err)
		}
	}
	return mydict
}

func TestLoaded(t *testing.T) {
	insertCases := []struct {
		str    string
		rtnErr error
	}{
		{"jS", nil},
		{"jS", nil},
		{"a", nil},
		{"a", nil},
	}
	lookupCases := []struct {
		str    string
		rtnStr string
		rtnErr string
	}{
		{"jS", "jS", ""},
		{"a", "a", ""},
		{"pizza", "", "entry not found"},
		{"b", "b", ""},
	}
	deleteCases := []struct {
		str    string
		rtnErr string
	}{
		{"jS", ""},
		{"jS", "entry not found, could not delete"},
		{"a", ""},
		{"a", "entry not found, could not delete"},
		{"pizza", "entry not found, could not delete"},
		{"b", ""},
		{"b", "entry not found, could not delete"},
	}
	mydict := loadTable()
	for i, item := range insertCases {
		_, err := mydict.Insert(DictEntry{Str: item.str})
		if item.rtnErr != err {
			t.Error("Unexpected Insert error [", err, "] for item[", i, "][", item.str, "]")
		}
	}
	for i, item := range lookupCases {
		resultEntry, err := mydict.Lookup(DictEntry{Str: item.str})
		if "" == item.rtnErr && err != nil {
			t.Error("Unexpected lookup error [", err, "] for item[", i, "][", item.str, "]")
		}
		if err != nil {
			if item.rtnErr != err.Error() {
				t.Error("Unexpected lookup error string [", err, "] for item[", i, "][", item.str, "]")
			}
			//fmt.Printf("HashDictError is: [%s]\n", err)
		}
		if item.rtnStr == "" {
			if resultEntry != nil {
				t.Errorf("100: expected nil return for item[%d]\n", i)
			}
		} else {
			if resultEntry.Str != item.rtnStr {
				t.Error("Expected lookup result", item.rtnStr, "but got", resultEntry.Str)
			}
		}
	}
	for i, item := range deleteCases {
		err := mydict.Delete(DictEntry{Str: item.str})
		if "" == item.rtnErr && err != nil {
			t.Error("Unexpected delete error [", err, "] for item[", i, "][", item.str, "]")
		}
	}
}

var dictCount int

func f(a DictEntry) {
	fmt.Printf("%d: %s\n", dictCount, a.Str)
	dictCount++
}

func TestMap(t *testing.T) {
	mydict := loadTable()
	mydict.Map(f)
}
