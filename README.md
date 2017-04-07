# utils 

This is a set of utilities / tools / library collections written in go. 

BackupCopy (basically normal copy of sourse dir to backup dir)

-	Strait copy of files to a back up location
o	Config file (not yet implemented)
	with list of directories to copy from
	directory to copy to
o	–n – copy new, never overwrite, only copy to a new location
o	–f – force copy, overwrite if needed
o	If neither –n or –f warn before over writing (currently (yes/no) should this include a upgrade to force all future copies or quit?)
o	Check of return types with size may be redundant if os.copy() returns an error if entire file is not copied. I need to check the copy() spec and update code if necessary

-	Analysis
o	Builds an analysis of all the files in a directory tree (multiple trees? Config file?)
	Produce an MD5 hash for each file and place in hashDictionary. All files with same MD5 will be linked to their single MD5 hash entry. Possibly a second hashDictionary based on file name.
	Possibly a second hashDictionary based on file name.
	Possibly a third hashDictionary based on file size.
	Each file object should have a structure with its size, hash, name, path.
	Each hashDictionary item will include a list of all files (structures) that match the entry. So for example: As with the hash case above, all the files of a given size will be linked to the same size entry in the hashdictionary of size. Likewise all the files of a given name can be found by looking up the name in the hashdictionary for names and its associated linked files.
	Dump all analysis into a file where it can be recovered/reloaded
	Filter flags to only print files with a single copy, only print files with multiple copies, flag for prime factor (md5 hash, name, ?)
	


-	Merging of two directories into one (or two backups into one)
o	Keep only one file when file is the same and they would copy into the same location
o	Analysis the difference when multiple files that would be placed into the same location but are not the same
	Real-time and off line

-	Restore : move file or the whole tree back from a backup location
o	All
o	Sub directory or file

- HashDictionary
o	This is ment to be a generic container.
o	Currently has string key (want to generalize this as needed but not sure how to best do this yet)
o	The value is not required but when desire some extra effort is required

To use HashDictionary you must:
- define the size of the hash table (should be a prime)
- define a hash function
- define an equal function
- Pass the above to Create() which will return the dictionary

If you want to add a value / payload to go along with the key the 
current method is to create a struct with the information you want and two pointer type conversion fuctions using the unsafe package. An example can be see in the analysis tool. 

