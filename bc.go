package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "net"
    "strconv"
    "strings"
)

const (
    TCP  = "tcp"
    IP = "localhost"
)
var nodes map[int]bool = make(map[int]bool)

func host(ip string, port int) string {
    return fmt.Sprintf("%s:%d", ip, port)
}
func getMsg(cn net.Conn) string {
    r := bufio.NewReader(cn)
    msg, _ := r.ReadString('\n')
    return strings.TrimSpace(msg)
}
func sendNoReply(toIp int, msg string) {
    cn, _ := net.Dial(TCP, host(IP, toIp))
    defer cn.Close()
    fmt.Fprintln(cn, msg)
}
func send(toIp int, msg string) string {
    cn, _ := net.Dial(TCP, host(IP, toIp))
    defer cn.Close()
    fmt.Fprintln(cn, msg)
    return getMsg(cn)
}
func ipAdder(chId <-chan int) {
    for {
        nodes[<-chId] = true
        fmt.Println(nodes)
    }
}

func servAdder(id int, chId chan<- int) {
    ln, _ := net.Listen(TCP, host(IP, id + 1))
    defer ln.Close()
    for {
        cn, _ := ln.Accept()
        go func(cn net.Conn) {
            val, _ := strconv.Atoi(getMsg(cn))
            chId<- val
            cn.Close()
        }(cn)
    }
}
func cliAdder(newId int, nodes map[int]bool) {
    msg := fmt.Sprintf("%d", newId)
    for target := range nodes {
       go sendNoReply(target + 1, msg)
    }
}
func servRegister(id int, chId, end chan<- int) {
    ln, _ := net.Listen(TCP, host(IP, id + 2))
    defer ln.Close()
    for {
        cn, _ := ln.Accept()
        go func(cn net.Conn) {
            newId, _ := strconv.Atoi(getMsg(cn))
            cliAdder(newId, nodes)
            buf, _ := json.Marshal(nodes)
            fmt.Fprintln(cn, string(buf))
            chId<- newId
            cn.Close()
        }(cn)
    }
    end<- 0
}
func cliRegister(id, targetId int, chId chan<- int) {
    resp := send(targetId + 2, fmt.Sprintf("%d", id))
    var slc map[int]bool
    _ = json.Unmarshal([]byte(resp), &slc)
    for newId := range slc {
        chId<- newId
    }
    // fmt.Println("Recibido: ", nodes)
}

func main() {
    func main() {
    chId := make(chan int)
    end := make(chan int)
    go ipAdder(chId)
    id := 0
    fmt.Print("Ingresa tu Port: ")
    fmt.Scanf("%d\n", &id)
    go servAdder(id, chId)
	go servRegister(id, chId, end)
	
	friend := 0
    fmt.Print("Ingresa tu friend Port: ")
	fmt.Scanf("%d\n", &friend)

    if id != friend {
        chId<- friend
        cliRegister(id, friend, chId)
    }
    <-end
}
}
