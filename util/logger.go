package util

import (
	"log"
	"os"
)

var (
	Logger *log.Logger
)

func init() {
	Logger = log.New(os.Stdout, "[Go-RealWorld] ", log.LstdFlags|log.Lshortfile)
}
