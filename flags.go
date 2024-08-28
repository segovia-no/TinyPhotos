package main

import (
	"flag"
)

type Flags struct {
	file           string
	bulkFromFolder string
	maxRoutines    int
}

func (f *Flags) parseFlags() {
	file := flag.String("file", "", "Compresses a single file providing the relative or absolute filepath")
	bulkFromFolder := flag.String("bulkfromfolder", "", "Compresses all the files in a folder providing the relative or absolute path to the folder")
	maxRoutines := flag.Int("maxroutines", 1, "Max number of routines to run concurrently")

	flag.Parse()

	f.file = *file
	f.bulkFromFolder = *bulkFromFolder
	f.maxRoutines = *maxRoutines
}
