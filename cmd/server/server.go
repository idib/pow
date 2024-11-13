package main

import (
	"context"
	"fmt"
	"log"
	"net"
)

type randomizerIntner interface {
	Intn(int) int
}

type RandomAccessResource interface {
	getRandom() string
}

type POWTask interface {
	CheckSolve(string) bool
	Present() string
}

type POWGenerator interface {
	MakeTask() POWTask
}

type Server struct {
	powGenerator        POWGenerator
	RandomQuoteProvider RandomAccessResource
	Port                int
}

func NewServer(randomQuoteProvider RandomAccessResource, port int, powGenerator POWGenerator) *Server {
	return &Server{
		RandomQuoteProvider: randomQuoteProvider,
		Port:                port,
		powGenerator:        powGenerator,
	}
}

func (s Server) check() {
	if s.Port <= 0 && s.Port > 65535 {
		log.Fatalf("Port must be in range 1-65535, but %d given", s.Port)
	}
	if s.RandomQuoteProvider == nil {
		log.Fatal("RandomQuoteProvider must be not nil")
	}
	if s.powGenerator == nil {
		log.Fatal("powGenerator must be not nil")
	}
}

func (s Server) run() {
	s.check()

	tcpServer := net.ListenConfig{
		Control:         nil,
		KeepAlive:       0,
		KeepAliveConfig: net.KeepAliveConfig{},
	}

	listener, err := tcpServer.Listen(context.Background(), "tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Fatal("critic error:", err)
	}
	defer listener.Close()

	log.Println("Server started on port ", s.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("connection error:", err)
			continue
		}

		// For now I have decided not to create a connection manager or connection pool.
		go s.handleConnection(conn)
	}
}

func (s Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("New connection from", conn.RemoteAddr())
	defer log.Println("Connection closed from", conn.RemoteAddr())

	// Here you can get the client's address and try to restore the previous session. I did not implement this, but left the possibility.
	csmp := NewCSMP(conn)
	err := csmp.Start()
	if err != nil {
		log.Println(err)
		return
	}

	task := s.powGenerator.MakeTask()

	err = csmp.SendTask(task.Present())
	if err != nil {
		log.Println(err)
		return
	}

	proofStr, err := csmp.WaitAnswer()
	if err != nil {
		log.Println(err)
		return
	}

	isOkAnswer := task.CheckSolve(proofStr)

	if isOkAnswer {
		err = csmp.SendResource(s.RandomQuoteProvider.getRandom())
	} else {
		err = csmp.SendWrongAnswer()
	}
	if err != nil {
		log.Println(err)
		return
	}
}
