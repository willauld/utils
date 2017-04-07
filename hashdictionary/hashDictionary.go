/*
I want to make this a generic hash dictionary library but I don't see
anyway to do this at the moment. So, for now, I'll just apply it to
strings.
*/

package hashdictionary

//func init() { fmt.Printf("NumCore: %d, Usable: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(0)) }

// ValuePayload is a place holder for your own data
type ValuePayload struct {
}

// DictEntry is entry type for an instantiation of HashDictionary
type DictEntry struct {
	Str   string
	Value *ValuePayload
}

type listEntry struct {
	next  *listEntry
	entry *DictEntry
}

// HashDictionary is a dictionary formed from a hashed array with linked list of coliding entries
type HashDictionary struct {
	table *[]*listEntry
	size  int
	hash  func(DictEntry) int // must hash modulo size
	equal func(DictEntry, DictEntry) bool
}

// Create constructs and initializes a HashDictionary
func Create(size int, hash func(DictEntry) int,
	equal func(DictEntry, DictEntry) bool) (dict *HashDictionary) {
	v := make([]*listEntry, size)
	dict = &HashDictionary{&v, size, hash, equal}
	return
}

// Insert adds an entry to the dictionary if it is not already present, returns the entry
func (dict *HashDictionary) Insert(insertEntry DictEntry) (ent *DictEntry, err error) {
	h := dict.hash(insertEntry)
	for entry := (*dict.table)[h]; entry != nil; entry = entry.next {
		if dict.equal(*entry.entry, insertEntry) {
			return entry.entry, nil
		}
	}
	e := &listEntry{(*dict.table)[h], &insertEntry}
	(*dict.table)[h] = e
	return e.entry, nil
}

// Lookup searches for a DictEntry item in the HashDictionary, returns
// the DictEntry or an error if not found
func (dict *HashDictionary) Lookup(lookupEntry DictEntry) (ent *DictEntry, err error) {
	h := dict.hash(lookupEntry)
	for entry := (*dict.table)[h]; entry != nil; entry = entry.next {
		if dict.equal(*entry.entry, lookupEntry) {
			return entry.entry, nil
		}
	}
	err = &hashDictError{"entry not found"}
	return
}

// Delete removes a DictEntry item from the HashDictionary, returns an error is not found
func (dict *HashDictionary) Delete(deleteEntry DictEntry) (err error) {
	h := dict.hash(deleteEntry)
	lptr := &(*dict.table)[h]
	for entry := (*dict.table)[h]; entry != nil; entry = entry.next {
		if dict.equal(*entry.entry, deleteEntry) {
			*lptr = entry.next
			return nil
		}
		lptr = &entry.next
	}
	err = &hashDictError{"entry not found, could not delete"}
	return
}

// Map call the supplied fuction f on each DictEntry item in the HashDictionary
func (dict *HashDictionary) Map(f func(DictEntry)) {
	for i := 0; i < dict.size; i++ {
		for entry := (*dict.table)[i]; entry != nil; entry = entry.next {
			f(*entry.entry)
		}
	}
}

type hashDictError struct {
	error string
}

func (e *hashDictError) Error() string {
	return e.error
}
