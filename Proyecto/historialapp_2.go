package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	//"time"
)

const ID = 1
const host = "localhost:8004"
const protocol = "tcp"
const webServerPort = ":5001"

var dirs = []string{"http://e96297a8.ngrok.io", "http://f2159382.ngrok.io"}
var localBlockchain Blockchain
var nodosBlockchain []Blockchain

type HistoriaClinica struct {
	Nombre        string `json:"nombre" binding:"required"`
	Edad          string `json:"edad" binding:"required"`
	Weight        string `json:"weight" binding:"required"`
	IsSmoker      bool   `json:"is_smoker" binding:"exists"`
	IsADrugUser   bool   `json:"is_a_drug_user" binding:"exists"`
	DrinksAlcohol bool   `json:"drinks_alcohol" binding:"exists"`
	FamilyStatus  string `json:"family_status" binding:"required"`
	Ocupation     string `json:"ocupation" binding:"required"`
	IsGenesis     bool   `json:"is_genesis"`
}

type Block struct {
	Pos       int             `json:"pos"`
	Data      HistoriaClinica `json:"data"`
	Timestamp string          `json:"timestamp"`
	Hash      string          `json:"hash"`
	PrevHash  string          `json:"prevhash"`
}

type Blockchain struct {
	Blocks []Block `json:"blocks"`
}

func GenesisBlock() Block {
	return CreateBlock(Block{}, HistoriaClinica{IsGenesis: true})
}

func NewBlockchain() Blockchain {
	return Blockchain{[]Block{GenesisBlock()}}
}

func (b Block) generateHash() Block {
	// obtener valoe de la Data
	bytes, _ := json.Marshal(b.Data)
	// concatenar el conjunto de datos
	data := string(b.Pos) + b.Timestamp + string(bytes) + b.PrevHash
	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
	return b
}

func CreateBlock(prevBlock Block, MedicalHistoryForm HistoriaClinica) Block {
	block := Block{}
	block.Pos = prevBlock.Pos + 1
	// se usa un tiempo fijo
	//block.Timestamp = time.Now().String()
	block.Timestamp = "00:00:00"
	block.Data = MedicalHistoryForm
	block.PrevHash = prevBlock.Hash
	block = block.generateHash()

	return block
}

func (bc Blockchain) AddBlock(data HistoriaClinica) Blockchain {
	// obtener bloque anterior
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	// crear nuevo bloque
	block := CreateBlock(prevBlock, data)
	//  validar integridad datos
	fmt.Println("Agrega bloque a la cadena de bloques local")
	if validBlock(block, prevBlock) {
		fmt.Println("Bloque valido")
		bc.Blocks = append(bc.Blocks, block)
		// retorna la cadena de bloques
		// se puede también usar punteros para no retornar la cadena de bloques
		return bc
	}
	return bc
}

func validBlock(block, prevBlock Block) bool {
	// Confirmar los hashes
	if prevBlock.Hash != block.PrevHash {
		return false
	}
	// confirmar el hash del bloque es válido
	if !block.validateHash(block.Hash) {
		return false
	}
	// Revisar la posicion para confirmar que ha sido incrementado
	if prevBlock.Pos+1 != block.Pos {
		return false
	}
	return true
}

func (b Block) validateHash(hash string) bool {
	b = b.generateHash()
	if b.Hash != hash {
		return false
	}
	return true
}

type Command struct {
	Origen int        `json:"origen"`
	Bl     Blockchain `json:"bl"`
	Tipo   int        `json:"tipo"`
}

func receiveAll(bl chan Command) {
	ln, _ := net.Listen(protocol, host)
	fmt.Printf("Escuchando en el puerto: %s\n", host)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go receive(conn, bl)
	}
}

func receive(conn net.Conn, bl chan Command) {
	alldirs := []string{"http://e96297a8.ngrok.io", "http://c0437b82.ngrok.io", "http://f2159382.ngrok.io"}

	defer conn.Close()
	r := bufio.NewReader(conn)
	msg, _ := r.ReadString('\n')
	fmt.Printf("Comando recibido: %s\n", msg)
	var c Command
	ba := []byte(msg[:len(msg)-1])
	json.Unmarshal(ba, &c)
	fmt.Printf("Se obtuve del origen %d el Blockchain\n", c.Origen)
	fmt.Println(c.Bl)

	if c.Tipo == 0 {
		localBlockchain = c.Bl
		command := &Command{Origen: ID, Bl: localBlockchain, Tipo: 1}
		out, _ := json.Marshal(command)
		jsonBlockchain := string(out)
		conn2, _ := net.Dial(protocol, alldirs[c.Origen])
		defer conn2.Close()
		fmt.Fprintf(conn2, "%s\n", jsonBlockchain)
	} else {
		fmt.Println("Respuesta de mensaje")
		bl <- Command{c.Origen, c.Bl, c.Tipo}
		fmt.Println("Comando recibido al canal")
	}

	//	time.Sleep(1 * time.Second)
	// responde al nodo con el blockchain

	// no asigna la cadena de bloques recibidos directamente a la cadena de bloques local
	// se debe hacer consenso para asignar
	//localBlockchain = c.Bl
	// nodosBlockchain[c.Origen] = c.Bl
	// fmt.Println(nodosBlockchain[c.Origen])
	//consenso()
}

func sendAll() {
	command := &Command{Origen: ID, Bl: localBlockchain, Tipo: 0}
	out, _ := json.Marshal(command)
	jsonBlockchain := string(out)
	fmt.Println(jsonBlockchain)
	for _, dir := range dirs {
		send(dir, jsonBlockchain)
	}
}

func send(remote, msg string) {
	conn, _ := net.Dial(protocol, remote)
	defer conn.Close()
	fmt.Fprintf(conn, "%s\n", msg)
}

func prepararConsenso() {
	n := len(dirs)
	fmt.Printf("Número de nodos: %d\n", n)
	nodosBlockchain = make([]Blockchain, n+1)

	for i := 0; i < n; i++ {
		nodosBlockchain[i] = Blockchain{}
	}
	fmt.Println(nodosBlockchain)
}

func consenso(c chan Command) {
	n := len(dirs)
	hashes := make([]string, n)
	//lastBlockList := make([]Block, n)
	blockchainNodos := make([]Blockchain, n+1)

	for j := 0; j < n; j++ {
		blockchainNodos[j] = Blockchain{}
	}

	nodosBlockchain[0] = localBlockchain

	//cont := 0
	for i := 0; i < n; i++ {
		blockchainReceived := <-c
		nodosBlockchain[blockchainReceived.Origen] = blockchainReceived.Bl
		bloques := blockchainReceived.Bl.Blocks
		ultimoBloque := bloques[len(bloques)-1]
		hashes[i] = ultimoBloque.Hash
		/*if len(bloques) != len(localBlockchain.Blocks) {
			cont++
		}
		lastBlockList[i] = ultimoBloque*/
		fmt.Println("Bloque recibido")
		fmt.Println(nodosBlockchain)
	}

	maxHashBloque := maxHash(hashes)
	fmt.Println("Hash del bloque de la mayoria: ")
	fmt.Println(maxHashBloque)

	for i := 0; i < n+1; i++ {
		fmt.Println("Bloque guardado:")
		fmt.Println(nodosBlockchain[i])
		bloques := nodosBlockchain[i].Blocks
		if len(bloques) == 0 {
			continue
		}

		fmt.Println(len(bloques))
		if bloques[len(bloques)-1].Hash == maxHashBloque {
			localBlockchain = nodosBlockchain[i]
			fmt.Println("Se asignó un nuevo Blockchain")
		}
	}

	//cont2 := 0

	/*for i := 0; i < n; i++ {
		if lastBlockList[i].Hash != maxHashBloque {
			cont2++
		}

	}

	if cont2 > 0 {
		fmt.Println("Problema de Hash. Actualizando hashes..")
		bloques := localBlockchain.Blocks
		ultimoBloque := bloques[len(bloques)-1]
		//ultimoBloque.Hash = maxHashBloque
		localBlockchain = ultimoBloque
	}*/

	// ultimoBloque := bloques[len(bloques) - 1]
	// 	hashes[i] = ultimoBloque.Hash
	// maxHashBloque := maxHash(hashes)

	//fmt.Printf("Hash de la mayoría: %s", maxHashBloque)

	// busca un bloque con el hash de la mayoría y lo asigna

	// if ultimoBloque.Hash == maxHashBloque {
	// 	localBlockchain = nodosBlockchain[n-1]
	// }

}

func maxHash(hashes []string) string {
	m := make(map[string]int)

	for _, h := range hashes {
		if m[h] == 0 {
			m[h] = 1
		} else {
			m[h] = m[h] + 1
		}
	}

	maxH := "none"
	maxHContador := 0

	for key, value := range m {
		if value > maxHContador {
			maxH = key
			maxHContador = value
		}
	}
	fmt.Println(maxHContador)
	return maxH
}

func main() {
	prepararConsenso()
	fmt.Println("Iniciado")
	localBlockchain = NewBlockchain()
	bl := make(chan Command)
	go receiveAll(bl)

	router := gin.Default()

	router.Use(cors.Default())

	/*historias := []HistoriaClinica {
		HistoriaClinica{
			Nombre: "Iván",
			Edad: "22",
			Weight: "65",
			IsSmoker: false,
			IsADrugUser: false,
			DrinksAlcohol: false,
			FamilyStatus: "Soltero",
			Ocupation: "Estudiante",
			IsGenesis: false,
		},
		HistoriaClinica{
			Nombre: "Persona",
			Edad: "30",
			Weight: "70",
			IsSmoker: true,
			IsADrugUser: false,
			DrinksAlcohol: true,
			FamilyStatus: "Casado",
			Ocupation: "A",
			IsGenesis: false,
		},
	}*/

	// Lista las historias clínicas
	router.GET("/historiasclinicas", func(c *gin.Context) {
		/*c.JSON(http.StatusOK, gin.H{
			"message": "historias clínicas",
		})*/
		go consenso(bl)
		go sendAll()

		bloques := localBlockchain.Blocks

		tamanio := len(bloques)
		if tamanio > 0 {
			historias := make([]HistoriaClinica, 0)

			for _, bloque := range bloques {
				historias = append(historias, bloque.Data)
			}

			c.JSON(http.StatusOK, historias)
		}
	})

	// Registra una historia clínica
	router.POST("/historiasclinicas", func(c *gin.Context) {
		var historia HistoriaClinica
		err := c.BindJSON(&historia)

		if err != nil {
			fmt.Println(err)
		}
		//nombre := c.PostForm("nombre")
		//edad := c.PostForm("edad")
		//fmt.Printf("nombre %s; edad %s", nombre, edad)
		localBlockchain = localBlockchain.AddBlock(historia)
		//go consenso(bl)

		go consenso(bl)
		go sendAll()

		fmt.Printf("Historia registrada: %v\n", historia)
		c.JSON(http.StatusOK, gin.H{
			"message": "historia clínica registrada",
		})
	})

	router.Run(webServerPort) // listen and serve on 0.0.0.0:8080
}
