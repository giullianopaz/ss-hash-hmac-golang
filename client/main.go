package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
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

func (m Message) String() string {
	return fmt.Sprintf(
		`{"name" : "%s", "text" : "%s", "hmac" : "%s"}`,
		m.Name, m.Text, m.HMAC,
	)
}

func main() {
	var (
		name      string
		ip        string
		port      int
		nMessages int
		alg       string
		// wg        sync.WaitGroup
	)

	flag.StringVar(&name, "name", "client_name", "Nome do Cliente")
	flag.StringVar(&ip, "ip", "localhost", "IP do Servidor")
	flag.IntVar(&port, "port", 8000, "Porta do Servidor")
	flag.IntVar(&nMessages, "n_messages", 100, "Quantidade de Mensagens a serem enviadas")
	flag.StringVar(&alg, "alg", "diffie-hellman", "Agoritmo para gerar a chave compartilhada")
	flag.Parse()

	// var dh DiffieHellman
	// dh.SetpModulusValue(47)
	// dh.SetgBaseValue(13)
	// dh.GeneratePrivateValue()
	// dh.GeneratePublicValue()
	// dh.GenerateSharedPrivateKey(31)

	// fmt.Printf("\nDiffieHellman: %+v\n", dh)
	// return

	fmt.Println("n_messages: ", nMessages)
	msg := &Message{
		Name: name,
		Text: "teste",
		HMAC: "hmac---fwfwlfkmlkwfmewlfml",
	}

	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8000")

	for i := 0; i < nMessages; i++ {
		// send to socket
		fmt.Fprintf(conn, msg.String()+"\n")

		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("[Server]: " + message)
	}
	fmt.Fprintf(conn, "END"+"\n")
	return
}
