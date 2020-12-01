package main

import (
	"log"
	"os"
	"runtime"
)

// Credits to https://twinnation.org/articles/35/how-to-add-colors-to-your-console-terminal-output-in-go for colourization.

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func init() {

	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}

func CheckError(err error) {

	if err != nil {
		log.Fatalf(Red + "[Error]" + Reset + ": " + err.Error())
		os.Exit(1)
	}

}

func RemoveFromBuffer(buffer []BufferPair, i int) []BufferPair {
	buffer[len(buffer)-1], buffer[i] = buffer[i], buffer[len(buffer)-1]
	return buffer[:len(buffer)-1]
}
