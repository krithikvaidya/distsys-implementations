package main

import (
	"log"
	"os"
)

func CheckError(err error) {

	if err != nil {
		log.Printf("<<Error>>: %s", err.Error())
		os.Exit(1)
	}

}

func RemoveFromBuffer(buffer []BufferPair, i int) []BufferPair {
	buffer[len(buffer)-1], buffer[i] = buffer[i], buffer[len(buffer)-1]
	return buffer[:len(buffer)-1]
}
