package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"time"
	"encoding/json"
)

var HOSTS = []string{}

const (
	PROTOCOL  = "tcp"
	LOCALHOST = "10.11.97.199:8000"
)

func send() {
	//host := HOSTS[rand.Intn(len(HOSTS))]
	con, _ := net.Dial(PROTOCOL, "10.11.98.201:8000")
	//fmt.Println(host)
	defer con.Close()
	fmt.Fprintln(con, LOCALHOST)
}

func recieve(con net.Conn) {
	defer con.Close()
	r := bufio.NewReader(con)
	msg, _ := r.ReadString('\n')
	_ = json.Unmarshal([]byte(msg),HOSTS)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	go send()

	ln, _ := net.Listen(PROTOCOL, LOCALHOST)
	defer ln.Close()

	for {
		con2, _ := ln.Accept()
		go recieve(con2)
	}
}