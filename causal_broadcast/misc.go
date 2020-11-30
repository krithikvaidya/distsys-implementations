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
