package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

// CEILINGVALUE : Valor máximo para ser usado como `pModulusValue`
var CEILINGVALUE int = 50

// RAND : Reconfigura Seed
var RAND *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Message struct to get message from JSON decoder
type Message struct {
	Name string `json:"name"`
	Text string `json:"text"`
	HMAC string `json:"hmac"`
}

// Check verifica a existência de erros
func Check(err error, errMessage string) {
	if err != nil {
		fmt.Println(errMessage, " -> ", err.Error())
		os.Exit(1)
	}
}

func main() {
	fmt.Println("[INFO] Starting Server...")

	var (
		name string
		port int
		alg  string
		// wg        sync.WaitGroup
	)

	flag.StringVar(&name, "name", "client_name", "Nome do Cliente")
	flag.IntVar(&port, "port", 8000, "Porta do Servidor")
	flag.StringVar(&alg, "alg", "diffie-hellman", "Agoritmo para gerar a chave compartilhada")
	flag.Parse()

	// var dh DiffieHellman
	// dh.GeneratepModulusValue(CEILINGVALUE)
	// dh.GenerategBaseValue()
	// dh.GeneratePrivateValue()
	// dh.GeneratePublicValue()
	// dh.GenerateSharedPrivateKey(10)

	// fmt.Printf("\nDiffieHellman: %+v\n", dh)
	// return

	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, err := net.Listen("tcp", ":8000")
	Check(err, "Erro ao abrir a conexão!!")

	// accept connection on port
	conn, _ := ln.Accept()

	// run loop forever (or until ctrl-c)
	defer ln.Close()
	for {
		// will listen for message to process ending in newline (\n)
		msg, _ := bufio.NewReader(conn).ReadString('\n')

		if msg == "END\n" {
			ln.Close()
			return
		}
		var message Message
		json.Unmarshal([]byte(msg), &message)
		// output message received
		fmt.Printf("\nMessage Received: %+v", message)
		// send new string back to client
		conn.Write([]byte("OK" + "\n"))
	}
}
