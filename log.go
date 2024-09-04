package main

import (
	"io"
	"log"
	"os"
	"time"
)


func setupLogToFile() *os.File {
	logFile, err := os.OpenFile(getLogFilename(), os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	return logFile
}

func getLogFilename() string {
	t := time.Now()
	return t.Format("20060102150405") + ".log"
}
