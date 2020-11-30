package main

import (
	"bufio"
	"flag"
	"log"
	"math/rand"
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

	log.Printf("Causal Broadcast Simulator\n")

	// bind

	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanLines)

	log.Printf("Press enter when all processes are online.")

	proceed := sc.Text()

}
