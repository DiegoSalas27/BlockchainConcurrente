package main

import(
    "fmt"
    "net"
    "bufio"
    "encoding/json"
)

const MYIP = "10.11.98.205:8000"

type Ledger struct{
    Nombre string
}

func CreateLedger () *Ledger{
    
    return &Ledger{"Juan"}
} 

func handle(con net.Conn) {
	defer con.Close()
	r := bufio.NewReader(con)
    msg, _ := r.ReadString('\n')
// todo verificar hola
    l:=CreateLedger()
    buf, _ := json.Marshal(*l)
    fmt.Fprintln(con, string(buf))
    fmt.Println("recibido : ", msg)
}

func main() {
    ln, _ := net.Listen("tcp", MYIP)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go handle(con)

	}

}
