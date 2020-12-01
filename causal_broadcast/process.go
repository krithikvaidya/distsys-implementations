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

type BufferPair struct {
	Pid   int
	Clock []int
}

type Vector_Clock struct {
	n_proc      int
	PID         int
	Causal_time []int
	Buffer      []BufferPair
	ClockMutex  sync.RWMutex
}

func InitializeClock(n_process, pid int) *Vector_Clock {
	return &Vector_Clock{
		Causal_time: make([]int, n_proc),
		Buffer:      make([]BufferPair, 0, 1000), // 2d slice of len 0 and capacity 1000
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

		var i int
		for i = 255; i >= 0; i-- {

			if msg[i] != 0 {
				break
			}

		}

		msg = msg[:i+1]
		// log.Printf("Message is %v ok\n", msg)
		var rcvd_msg map[string][]int

		json.Unmarshal(msg, &rcvd_msg)

		// log.Printf("Message rcvd from PID: %v with clock %v\n", rcvd_msg["pid"], rcvd_msg["clock"])
		// log.Printf("\nrcvd msg is %v uf type %T\n", rcvd_msg, rcvd_msg)

		rcvd_clock, _ := rcvd_msg["clock"]
		// log.Printf("\n\nOK? : %v\n\n", ok)
		// if !ok {
		// 	log.Printf("got data of type %T but wanted []int", rcvd_msg["clock"])
		// }

		// log.Printf("RCVD clock %v and size %v\n", rcvd_clock, len(rcvd_clock))
		rcvd_pid1, _ := rcvd_msg["pid"]
		// log.Printf("\n\nOK? : %v, PID %v \n\n", ok, rcvd_pid1)
		rcvd_pid := rcvd_pid1[0]

		vclock.ClockMutex.Lock() // not RLock and then Lock

		immediate_deliver := true

		for i := 0; i < len(vclock.Causal_time); i++ {

			if i == rcvd_pid {

				if vclock.Causal_time[i] != rcvd_clock[i]-1 {

					// some more message(s) need to be delivered from the same sender proces
					// delivering this message.
					immediate_deliver = false
					break

				}
			} else {
				if vclock.Causal_time[i] < rcvd_clock[i] {

					// some more message(s) need to be delivered from other sender
					// process(es)
					immediate_deliver = false
					break

				}
			}
		}

		if immediate_deliver {

			vclock.Causal_time[rcvd_pid]++

			log.Printf(Green+"[Delivery Success]"+Reset+": Immediately delivered message from PID: %v with clock %v.\n"+Purple+"[Clock Value]"+Reset+": Current value of clock is %v\n", rcvd_pid, rcvd_clock, vclock.Causal_time)

			// deliver other buffered messages ready for delivery

			for i := 0; i < len(vclock.Buffer); i++ {

				diff := 0
				for j := 0; j < len(vclock.Buffer[i].Clock); j++ {

					if vclock.Buffer[i].Clock[j] > vclock.Causal_time[j] {

						diff += vclock.Buffer[i].Clock[j] - vclock.Causal_time[j]

					}

				}

				if diff <= 1 { // deliver this message

					for j := 0; j < len(vclock.Buffer[i].Clock); j++ {

						if vclock.Buffer[i].Clock[j] > vclock.Causal_time[j] {

							vclock.Causal_time[j]++
							break

						}

					}

					log.Printf(Green+"[Delivery Success]"+Reset+": Delivered buffered message from PID: %v with clock %v.\n"+Purple+"[Clock Value]"+Reset+": Current value of clock is %v\n", vclock.Buffer[i].Pid, vclock.Buffer[i].Clock, vclock.Causal_time)

					// remove from buffer
					vclock.Buffer = RemoveFromBuffer(vclock.Buffer, i)
					i--
				}

			}

		} else {

			// buffer it
			vclock.Buffer = append(vclock.Buffer, BufferPair{
				Pid:   rcvd_pid,
				Clock: rcvd_clock,
			})

			log.Printf(Blue+"[Buffered Message]"+Reset+": Buffered message from PID: %v with clock %v\n.", rcvd_pid, rcvd_clock)

		}

		vclock.ClockMutex.Unlock()

	}

}

func (vclock *Vector_Clock) CreateMessageListeners(listener *net.TCPListener) {

	for {

		conn, err := listener.Accept()

		if err != nil {
			// todo
			os.Exit(1)
		}

		log.Printf(Cyan + "[Info]" + Reset + ":" + fmt.Sprintf("Accepted an incoming connection request from [%s].", conn.RemoteAddr()))

		go vclock.ListenForMessages(conn)

	}

}

func (vclock *Vector_Clock) SendMessage(conn net.Conn, to_send []byte) {

	max := 15
	min := 5

	seconds := rand.Intn(max-min) + min

	time.Sleep(time.Duration(seconds) * time.Second)

	conn.Write(to_send)

}

func (vclock *Vector_Clock) CreateMessages(connxns []net.Conn) {

	max := 15 // max time to wait before sending message, in sec
	min := 5

	// Create messages to be broadcast at random times
	for {

		seconds := rand.Intn(max-min) + min
		log.Printf(Cyan+"[Info]"+Reset+": Waiting for %v seconds", seconds)
		time.Sleep(time.Duration(seconds) * time.Second)

		vclock.ClockMutex.Lock()

		vclock.Causal_time[vclock.PID] = vclock.Causal_time[vclock.PID] + 1 // increment self clock to record send event

		to_send := make(map[string][]int)
		to_send["pid"] = make([]int, 1, 1)
		to_send["pid"][0] = vclock.PID
		to_send["clock"] = vclock.Causal_time

		// Marshal the map into a slice of bytes.
		to_send_bytes, err := json.Marshal(to_send)
		CheckError(err)

		size := 256 - len(to_send_bytes)
		padding_bytes := make([]byte, size)

		to_send_bytes = append(to_send_bytes, padding_bytes...)
		to_send_bytes = to_send_bytes[:256] // just to ensure capacity is 256

		// to_send_str := string(to_send_bytes)
		// to_send_str = fmt.Sprintf("%-256v", to_send_str)

		log.Printf(Yellow+"[Send]"+Reset+": Broadcasting vector clock with values %v \n"+Purple+"[Clock Value]"+Reset+": Current Clock Value %v\n", string(to_send_bytes), vclock.Causal_time)

		vclock.ClockMutex.Unlock()

		for i := 0; i < vclock.n_proc-1; i++ {

			go vclock.SendMessage(connxns[i], to_send_bytes)

		}

		// log.Printf(Green+"[Success]"+Reset+": Successfully broadcasted message with clock %v\n", to_send["clock"])

	}

}
