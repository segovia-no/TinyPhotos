package main

import (
	"flag"
)

type Flags struct {
	apikey         string
	file           string
	bulkFromFolder string
	maxRoutines    int
	log            bool
	maxRetries     int
}

func (f *Flags) parseFlags() {
	apikey := flag.String("apikey", "", "You can pass the Tinify API key here instead of the .env file")
	file := flag.String("file", "", "Compresses a single file providing the relative or absolute filepath")
	bulkFromFolder := flag.String("bulkfromfolder", "", "Compresses all the files in a folder providing the relative or absolute path to the folder")
	maxRoutines := flag.Int("maxroutines", 4, "Max number of routines to run concurrently")
	log := flag.Bool("log", false, "Writes the generated log into a file")
	maxRetries := flag.Int("maxretries", 2, "Max number of retries if a request to the Tinify API fails per file")

	flag.Parse()

	f.apikey = *apikey
	f.file = *file
	f.bulkFromFolder = *bulkFromFolder
	f.maxRoutines = *maxRoutines
	f.log = *log
	f.maxRetries = *maxRetries
}
