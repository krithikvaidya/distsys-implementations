package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

var n_proc int

func init() {

	flag.IntVar(&n_proc, "n", 3, "number of processes (default 3)")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

}

func main() {

	fmt.Println()
	log.Println("Causal Broadcast Simulator\n")

	log.Println("Enter the port number the process should bind to: ")
	var port string
	fmt.Scanf("%s", &port)
	port = ":" + port

	// listen
	_, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}

	log.Println("Successfully bound to ", port)

	log.Println("Press enter when all processes are online.")

	var input rune

	fmt.Scanf("%c", &input)

}
