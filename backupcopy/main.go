/*
** original code based on: https://www.socketloop.com/tutorials/golang-copy-directory-including-sub-directories-files
** with pieces from: http://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
 */

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
)

var version = struct {
	major int
	minor int
}{0, 1}
var (
	force     bool
	newTarget bool
)

func getTerminalAnswer() bool {
	var str string
	for {
		fmt.Printf("(yes or no):")
		fmt.Scanln(&str)
		if str == "yes" || str == "y" {
			return true
		} else if str == "no" || str == "n" {
			return false
		}
	}
}

func dirExistsContinue(path string) bool {
	// ask user in terminal is over write is OK
	fmt.Printf("Dir: %s exists, continue?\n", path)
	// TODO: read in the terminal input
	return getTerminalAnswer()
}
func overWriteOK(path string) bool {
	// ask user in terminal is over write is OK
	fmt.Printf("File: %s exists, overwrite?\n", path)
	// TODO: read in the terminal input
	return getTerminalAnswer()
}

func copyFile(source string, dest string, size int64) (err error) {
	sfile, err := os.Stat(source)
	if err != nil {
		return
	}
	if !sfile.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		//return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
		fmt.Printf("CopyFile: non-regular source file %s (%q)",
			sfile.Name(), sfile.Mode().String())
		return
	}
	dfile, err := os.Stat(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			// OK is file does not exist but fail(return) if any other error
			return //TODO: ?????????????????
		}
	} else {
		if !(dfile.Mode().IsRegular()) {
			//return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
			fmt.Printf("CopyFile: non-regular destination file %s (%q)",
				dfile.Name(), dfile.Mode().String())
			return
		}
	}
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	var destFile *os.File
	if force || newTarget {
		// newTarget means the whole destination tree is new
		// so no checks are needed at this level.
		destFile, err = os.Create(dest)
		if err != nil {
			return err
		}

	} else {
		destFile, err = os.OpenFile(dest, os.O_EXCL|os.O_CREATE /*getFromSrc*/, 0666)
		if err != nil {
			if os.IsExist(err) {
				if !overWriteOK(dest) {
					fmt.Printf("File %s Exists, will not over write\n",
						dest)
					return nil
				}
				destFile, err = os.Create(dest)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	defer func() {
		cerr := destFile.Close()
		if err == nil {
			err = cerr
		}
	}()

	var bytes int64
	bytes, err = io.Copy(destFile, sourceFile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode()) //TODO: could this be done in the OPENFILE above?
		}

	}
	if bytes != size {
		fmt.Printf("Only %d of %d bytes copied for %s\n",
			bytes, size, source)
	}

	return
}

func copyDir(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir if nessasary
	_, err = os.Open(dest)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dest, sourceinfo.Mode())
		if err != nil {
			return err
		}
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourceFileStr := source + "/" + obj.Name()

		destinationFileStr := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = copyDir(sourceFileStr, destinationFileStr)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = copyFile(sourceFileStr, destinationFileStr, obj.Size())
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func main() {
	versionPtr := pflag.Bool("version", false, "program version")
	newTargetPtr := pflag.Bool("new", false, "new destination only")
	forcePtr := pflag.Bool("force", false, "force the copy, no questions")
	srcPtr := pflag.String("src", "", "source directory")
	destPtr := pflag.String("dest", "", "destination directory")
	pflag.Parse()
	fmt.Println("input:", *srcPtr)
	fmt.Println("input:", *destPtr)
	fmt.Println("tail:", pflag.Args())

	if *versionPtr == true {
		fmt.Printf("\t Version %d.%d", version.major, version.minor)
		os.Exit(0)
	}
	sourceDir := *srcPtr
	destDir := *destPtr
	force = *forcePtr
	newTarget = *newTargetPtr

	if sourceDir == "" {
		fmt.Printf("sourceDir is the empty string\n")
		sourceDir = "c:/home/auld/temp/backupCopyTest"
		destDir = "\\\\TOMATOHOST\\homeStore\\auld-THINK-backup\\dest1"
	}
	if destDir == "" {
		fmt.Printf("Must supply a source and destination\n")
		return
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

	// create the destination directory
	fmt.Println("Destination :" + destDir)

	if !*forcePtr {
		_, err = os.Open(destDir)
		if !os.IsNotExist(err) {
			if *newTargetPtr {
				fmt.Println("Destination directory already exists. Abort!")
				os.Exit(1)
			} else {
				if !dirExistsContinue(destDir) {
					os.Exit(0)
				}
			}
		}
	}

	err = copyDir(sourceDir, destDir)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Directory copied")
	}
}

/*
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var sourceDir, destDir string

func a(path string, f os.FileInfo, err error) error {
	fmt.Printf("f.Name: %s Size: %d Dir?: %t\n",
		f.Name(), f.Size(), f.IsDir())
	fmt.Println(path)
	return nil
}

func main() {
	sourceDir = ""
	destDir = ""
	flag.Parse()
	sourceDir = flag.Arg(0)
	destDir = flag.Arg(1)

	if sourceDir == "" {
		fmt.Printf("sourceDir is the empty string\n")
		sourceDir = "c:/home/auld/temp/backupCopyTest"
	}

	err := filepath.Walk(sourceDir, a)
	if err != nil {
		fmt.Printf("filepath.Walk failed %v\n", err)
		return
	}
}
*/
