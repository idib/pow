package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// CSMP - custom simple message protocol
type CSMP struct {
	state CSMPState
	conn  net.Conn
}

type CSMPState byte

const (
	initState CSMPState = iota
	start
	waitAnswer
	checkAnswer
	sendResource
	end
	abort
)

func NewCSMP(conn net.Conn) *CSMP {
	return &CSMP{state: initState, conn: conn}
}

func (csmp *CSMP) Start() error {
	// csmp.state == end If we have finished one communication, we can restart the protocol

	if csmp.state != initState && csmp.state != end {
		csmp.state = abort
		return errors.New("CSMP already started")
	}

	csmp.state = start

	return nil
}

func (csmp *CSMP) SendTask(taskVal string) error {
	if csmp.state != start {
		csmp.state = abort
		return errors.New("CSMP not started")
	}

	_, err := csmp.conn.Write([]byte(taskVal))
	if err != nil {
		csmp.state = abort
		return fmt.Errorf("error write to connection: %w", err)
	}
	csmp.state = waitAnswer
	return nil
}

func (csmp *CSMP) WaitAnswer() (answer string, err error) {
	if csmp.state != waitAnswer {
		csmp.state = abort
		return "", errors.New("CSMP in wrong state")
	}

	buffer := make([]byte, 1024)
	n, err := csmp.conn.Read(buffer)
	if err != nil {
		csmp.state = abort
		return "", fmt.Errorf("error read from connection: %w", err)
	}
	csmp.state = checkAnswer

	answer = strings.TrimSpace(string(buffer[:n]))
	return answer, nil
}

func (csmp *CSMP) SendResource(str string) error {
	if csmp.state != checkAnswer {
		return errors.New("CSMP not ready to send resource")
	}
	csmp.state = sendResource

	msg := fmt.Sprintf("recource:%s\n", str)
	_, err := csmp.conn.Write([]byte(msg))
	if err != nil {
		csmp.state = abort
		return fmt.Errorf("error write to connection: %w", err)
	}

	csmp.state = end
	return nil
}

func (csmp *CSMP) SendWrongAnswer() error {
	if csmp.state != checkAnswer {
		return errors.New("CSMP not ready to send resource")
	}
	_, err := csmp.conn.Write([]byte("wrong answer"))
	if err != nil {
		csmp.state = abort
		return fmt.Errorf("error write to connection: %w", err)
	}
	csmp.state = abort
	return nil
}

func (csmp *CSMP) ToAbort() {
	csmp.state = abort
}
