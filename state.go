package main

import (
	"fmt"
	"time"
)

type state struct {
	sendNum int
	recvNum int
	succNum int
	runTime int64
	qPSec   int
	startat int64
}

var State *state

const waitForClose = 10

func NewState() *state {
	return &state{}
}

func (s *state) Start() {
	s.startat = time.Now().UnixNano()
}

func (s *state) End() {
	tm := time.Now().UnixNano()
	s.runTime = (tm-s.startat)/1e9 - waitForClose +1
}

func (s *state) Show() {
	s.End()
	fmt.Printf("send requests:			%d\n", s.sendNum)
	fmt.Printf("receive responses:		%d\n", s.recvNum)
	fmt.Printf("responses success counts:	%d\n", s.succNum)
	fmt.Printf("time callapsed:			%d\n", s.runTime)
	var per int64
	if s.runTime == 0 {
		per = int64(s.sendNum)
	} else {
		per = int64(s.recvNum) / s.runTime
	}
	fmt.Printf("query per second:		%d\n", per)
}
