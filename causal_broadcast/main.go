package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"
)

var n_proc int

func init() {

	flag.IntVar(&n_proc, "n", 3, "number of processes (default 3)")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

}

func main() {

	fmt.Println("\nCausal Broadcast Simulator\n")

	fmt.Printf("Enter the process id: ")
	var pid int
	fmt.Scanf("%d", &pid)

	fmt.Printf("Enter the port number the process should bind to: ")
	var port string
	fmt.Scanf("%s", &port)
	port = ":" + port

	// listen
	_, err := net.Listen("tcp", port)
	CheckError(err)

	fmt.Println("Successfully bound to", port, "\n")

	fmt.Println("Press enter when all processes are online.")

	var input rune

	fmt.Scanf("%c", &input)

	fmt.Printf("Enter the port numbers of the other %v processes: \n", (n_proc - 1))

	ports := make([]string, n_proc-1)

	for i := 0; i < n_proc-1; i++ {
		fmt.Scan(&ports[i])
		ports[i] = ":" + ports[i]
	}

	connxns := make([]net.Conn, n_proc-1)

	for i := 0; i < n_proc-1; i++ {

		serverTcpAddr, err := net.ResolveTCPAddr("tcp", ports[i])
		CheckError(err)

		conn, err := net.DialTCP("tcp", nil, serverTcpAddr)
		CheckError(err)

		connxns[i] = conn
	}

	fmt.Println("Successfully connected to all processes.")

	clock := InitializeClock(n_proc, pid)

	clock.HandleMessageReception()
	clock.CreateAndSendMessages(ports)

}
