package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var debugFile *os.File
var logger *log.Logger

func init() {
	var err error
	filePath := filepath.Join(".", "debug.log")
	debugFile, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	logger = log.New(debugFile, "", log.LstdFlags|log.Lmicroseconds)
}

func Log(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logger.Println(message)
	debugFile.Sync() // Flush immediately
}

func Close() {
	if debugFile != nil {
		debugFile.Close()
	}
}
