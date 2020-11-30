package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
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

func (vclock *Vector_Clock) HandleMessageReception() {

}

func (vclock *Vector_Clock) CreateAndSendMessages(connxns []net.Conn) {

	max := 15 // max time to wait before sending message, in sec
	min := 5

	// Create messages to be broadcast at random times
	for {

		seconds := rand.Intn(max-min) + min
		log.Printf("Waiting for ", seconds, " seconds")
		time.Sleep(time.Duration(seconds) * time.Second)

		vclock.ClockMutex.Lock()

		vclock.Causal_time[pid] = vclock.Causal_time[pid] + 1 // increment self clock to record send event

		to_send := fmt.Sprintf("%d, ui", username, "broadcast", msg_str)

		for i = 0; i < vclock.n_proc-1; i++ {

		}

		vclock.ClockMutex.Unlock()

	}

}
