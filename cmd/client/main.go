package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func work(serverAdders string) error {
	log.Println("try connect to server:", serverAdders)
	conn, err := net.Dial("tcp", serverAdders)
	if err != nil {
		log.Println("connect error:", err)
		return err
	}
	defer conn.Close()
	log.Println("connected to server:", conn.RemoteAddr().String())

	message, _ := bufio.NewReader(conn).ReadString('\n')
	parts := strings.Split(strings.TrimSpace(message), ":")
	if len(parts) != 2 {
		log.Println("error: invalid task format")
		return err
	}
	log.Println("recive task:", parts[0], ":", parts[1])

	nonce, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Println("error: invalid nonce format")
		return err
	}
	difficulty, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Println("error: invalid difficulty format")
		return err
	}

	proof := -1
	var hashHex string
	target := strings.Repeat("0", difficulty)

	found := false
	startT := time.Now()
	for !found {
		proof++
		hashInput := fmt.Sprintf("%d%d", nonce, proof)
		hash := sha256.Sum256([]byte(hashInput))
		hashHex = hex.EncodeToString(hash[:])

		if strings.HasPrefix(hashHex, target) {
			found = true
		}
	}
	duration := time.Since(startT)

	log.Println("found proof:", proof)
	fmt.Fprintf(conn, "%d\n", proof)

	response, _ := bufio.NewReader(conn).ReadString('\n')
	log.Println(strings.TrimSpace(response))
	log.Println("attemts :", proof)
	log.Println("duration:", duration)
	return nil
}

func main() {
	serverAddress := flag.String("server", "localhost:3333", "The address of the server")
	flag.Parse()

	for i := 0; i < 5; i++ {
		if err := work(*serverAddress); err != nil {
			log.Println("error:", err)
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}
}
