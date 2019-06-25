package main
import (
	"fmt"
	"net"
	"bufio"
	"encoding/json"
)

type Ledger struct {
    Nombre string
}

var HOSTS = []string{"10.11.97.199:8000",
                     "10.11.98.201:8000",
                     "10.11.98.213:8000",
                     "10.11.98.211:8000",
                     "10.11.98.210:8000",
                     "10.11.98.215:8000",
                     "10.11.98.214:8000",
                     "10.11.98.226:8000",
                     "10.11.97.225:8000",
                     "10.11.97.218:8000",
                     "10.11.97.219:8000",
                     "10.11.98.207:8000",
                     "10.11.98.205:8000",
                     "10.11.98.229:8000"}

func sender(host string){
	con, _ := net.Dial("tcp", host)
	defer con.Close()
	r := bufio.NewReader(con)
	fmt.Fprint(con, "hola")
	ledg := &Ledger{}
	newmsg, _ := r.ReadString('\n')
	_ = json.Unmarshal([]byte(newmsg), ledg)
	fmt.Println(ledg)
}

func main(){
	for {
		fmt.Println("Press enter to send..")
		var input string
		fmt.Scanln(&input)
		sender("10.11.98.205:8000")
	}
}