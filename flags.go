package main

import (
	"flag"
)

type Flags struct {
	file           string
	bulkFromFolder string
}

func (f *Flags) parseFlags() {
	file := flag.String("file", "", "Compresses a single file providing the relative or absolute filepath")
	bulkFromFolder := flag.String("bulkfromfolder", "", "Compresses all the files in a folder providing the relative or absolute path to the folder")

	flag.Parse()

	f.file = *file
	f.bulkFromFolder = *bulkFromFolder
}
