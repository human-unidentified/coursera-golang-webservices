package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

func filterFiles(files []os.FileInfo, printFiles bool) []os.FileInfo {
	vsf := make([]os.FileInfo, 0)
	for _, v := range files {
		if !printFiles && !v.IsDir() {
			continue

		}
		vsf = append(vsf, v)
	}

	return vsf
}

type byName []os.FileInfo

func (s byName) Len() int {
	return len(s)
}
func (s byName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byName) Less(i, j int) bool {
	return s[i].Name() < s[j].Name()
}
func dirTreeProcessLevel(output io.Writer, path string, printFiles bool, level int, lastElementInParent bool, prevIndent string) {
	files, _ := ioutil.ReadDir(path)
	files = filterFiles(files, printFiles)
	sort.Sort(byName(files))

	const tab = "\t"

	newIndent := ""
	if level > 0 {
		newIndent = prevIndent + "│" + tab
	}

	if lastElementInParent {
		newIndent = prevIndent + tab
	}

	for idx, f := range files {
		lastElement := idx == len(files)-1
		fileName := f.Name()

		var prefix string
		if lastElement {
			prefix = "└───"
		} else {
			prefix = "├───"
		}

		if f.IsDir() == true {
			fmt.Fprintln(output, newIndent+prefix+fileName)
			var newpath = path + string(os.PathSeparator) + fileName

			dirTreeProcessLevel(output, newpath, printFiles, level+1, lastElement, newIndent)
		} else if printFiles {
			size := "empty"
			if fileSize := f.Size(); fileSize > 0 {
				size = strconv.FormatInt(f.Size(), 10) + "b"
			}

			fmt.Fprintln(output, newIndent+prefix+fileName+" ("+size+")")
		}
	}
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	dirTreeProcessLevel(output, path, printFiles, 0, false, "")
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
