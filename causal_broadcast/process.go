package main

import "sync"

type Vector_Clock struct {
	Causal_time []string
	ClockMutex  sync.RWMutex
}

func InitializeClock(n_proc int) {
	return &Vector_Clock{
		Causal_time: make([]int, n_proc),
	}
}

func (vclock *Vector_Clock) ListenForMessages() {

}

func (vclock *Vector_Clock) CreateAndSendMessages() {

}
