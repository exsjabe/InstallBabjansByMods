package main

import (
	"os"
	"strings"
	"sync"
)

// The information needed to extract a zipfile to it's destination in .minecraft.
type ExtractBundle struct {
	zipfile     string
	destination string
}

// A slice of ExtractBundle structs.
type ExtractRecord []ExtractBundle

// includesZipFile returns true if the record contains the given zipfile.
func (ea ExtractRecord) includesZipfile(zipfile string) bool {
	for _, v := range ea {
		if v.zipfile == zipfile {
			return true
		}
	}
	return false
}

// Extract extracts the record's zipfiles to their destination.
// src corresponds to the source directory, workDir corresponds to the temporary directory, and dest corresponds to the destination directory.
func (ea ExtractRecord) extract(src, workDir, dest string, res chan<- string) {
	wg := sync.WaitGroup{}
	for _, eb := range ea {
		wg.Add(1)
		go func(eb ExtractBundle) {
			defer wg.Done()
			eb.extractFile(src, workDir, dest, res)
		}(eb)
	}

	wg.Wait()
}

func (eb ExtractBundle) extractFile(src, workDir, dest string, res chan<- string) {

	err := Unzip(src+"\\"+eb.zipfile, workDir)

	CheckPanic(err)

	os.MkdirAll(dest+"\\"+eb.destination, 0755)

	unzippedDir := workDir + "\\" + strings.Split(eb.zipfile, ".zip")[0]
	wg := sync.WaitGroup{}

	files, err := os.ReadDir(unzippedDir)
	CheckPanic(err)
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			moveExtractedFile(file, unzippedDir, eb.destination, dest, res)
		}(file.Name())
	}
	wg.Wait()
}

// moveExtractedFiles moves the files from the workingdirectory to the targetFolder in the target directory.
func moveExtractedFile(file, workDir, targetFolder, targetDir string, res chan<- string) {

	if !PathExists(targetDir + "\\" + targetFolder + "\\" + file) {
		err := CopyFile(workDir+"\\"+file, targetDir+"\\"+targetFolder+"\\"+file)
		CheckPanic(err)
		res <- ("Successfully added " + file + " to " + targetFolder + ".")
	} else {
		res <- (file + " already exists in " + targetFolder + ".")
	}

}
