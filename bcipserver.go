package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	//"strings"
)
/*
type MessageType int32

const (
	NEWHOST MessageType = 0
	ADDHOST MessageType = 1
)

type RequestBody struct {
	message string
	messageType int32 
}
*/
const (
	PROTOCOL  = "tcp"
	LOCALHOST = "10.11.98.201:8000"
)

var HOSTS = []string{"10.11.98.201:8000"}

func handle(con net.Conn) {
	defer con.Close()
	r := bufio.NewReader(con)
	ip, _ := r.ReadString('\n')
	clientConn, _ := net.Dial("tcp", ip)
	defer clientConn.Close()
	/*request := &RequestBody{}
	body, _ := r.ReadString('\n')
	//msg = strings.TrimSpace(msg)
	_ = json.Unmarshal([]byte(body), request)

*/
	var flhost bool
	for _, element := range HOSTS {
		if element == ip{
			flhost=true
			break
		}
	}
	buf, _ := json.Marshal(HOSTS)
	fmt.Fprintln(clientConn, buf)

	if !flhost{
		HOSTS= append(HOSTS,string(ip))
	}
}

func main() {
	ln, _ := net.Listen(PROTOCOL, LOCALHOST)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go handle(con)
	}
}
