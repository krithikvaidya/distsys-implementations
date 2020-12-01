package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

type Vector_Clock struct {
	n_proc      int
	PID         int
	Causal_time []int
	ClockMutex  sync.RWMutex
}

func InitializeClock(n_process, pid int) *Vector_Clock {
	return &Vector_Clock{
		Causal_time: make([]int, n_proc),
		n_proc:      n_process,
		PID:         pid,
	}
}

func (vclock *Vector_Clock) ListenForMessages(conn net.Conn) {

	for {

		msg := make([]byte, 256)
		_, err := io.ReadFull(conn, msg) // read 256 bytes message

		if err != nil {
			// todo
			os.Exit(1)
		}

		msg_str := string(msg)

		var i int
		for i = 255; i >= 0; i-- {

			if msg_str[i] != ' ' {
				break
			}

		}

		msg = msg[:i+1]

		var rcvd_msg map[string]interface{}

		json.Unmarshal(msg, &rcvd_msg)

		log.Printf("Message rcvd from PID: %v with clock %v\n", rcvd_msg["pid"], rcvd_msg["clock"])

		vclock.ClockMutex.Lock()  // not RLock and then Lock

		immediate_deliver := true

		for i := 0; i < vclock.Causal_time.Size(); i++ {

			if i == rcvd_msg["pid"] {

				if vclock.Causal_time[i] != rcvd_msg["clock"][i] - 1 {

					// some more message(s) need to be delivered from the same sender proces
					// delivering this message.
					immediate_deliver = false
					break

				}

				else {
					if vclock.Causal_time[i] < rcvd_msg["clock"][i] {

						// some more message(s) need to be delivered from other sender
						// process(es)
						immediate_deliver = false
						break

					}
				}
			}
		}

		if immediate_deliver {
			vclock.Causal_time[rcvd_msg["pid"]]++
		}
		else {
			
		}

	}

}

func (vclock *Vector_Clock) CreateMessageListeners(listener *net.TCPListener) {

	for {

		conn, err := listener.Accept()

		if err != nil {
			// todo
			os.Exit(1)
		}

		log.Printf(fmt.Sprintf("Accepted an incoming connection request from [%s].", conn.RemoteAddr()))

		go vclock.ListenForMessages(conn)

	}

}

func (vclock *Vector_Clock) SendMessage(conn net.Conn, to_send string) {

	max := 15
	min := 5

	seconds := rand.Intn(max-min) + min
	log.Printf("Waiting for %v seconds before sending to process with conn %v\n", seconds, conn)
	time.Sleep(time.Duration(seconds) * time.Second)

	conn.Write([]byte(to_send))

}

func (vclock *Vector_Clock) CreateMessages(connxns []net.Conn) {

	max := 15 // max time to wait before sending message, in sec
	min := 5

	// Create messages to be broadcast at random times
	for {

		seconds := rand.Intn(max-min) + min
		log.Printf("Waiting for %v seconds", seconds)
		time.Sleep(time.Duration(seconds) * time.Second)

		vclock.ClockMutex.Lock()

		vclock.Causal_time[vclock.PID] = vclock.Causal_time[vclock.PID] + 1 // increment self clock to record send event

		to_send := make(map[string]interface{})
		to_send["pid"] = vclock.PID
		to_send["clock"] = vclock.Causal_time

		// Marshal the map into a slice of bytes.
		to_send_bytes, err := json.Marshal(to_send)
		to_send_str := string(to_send_bytes)
		to_send_str = fmt.Sprintf("%-256v", to_send_str)

		CheckError(err)

		vclock.ClockMutex.Unlock()

		log.Printf("\nBroadcasting vector clock with values %v \n", string(to_send_bytes))

		for i := 0; i < vclock.n_proc-1; i++ {

			go vclock.SendMessage(connxns[i], to_send_str)

		}

		log.Printf("Successfully broadcasted.\n")

	}

}
