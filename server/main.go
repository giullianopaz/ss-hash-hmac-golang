package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	// OKCOLOR : Para mostrar OK verde
	OKCOLOR = "\033[1;32m-> OK:  %+v\033[0m\n\n"
	// ERRORCOLOR : Para mostrar HMAC ERROR vermelho
	ERRORCOLOR = "\033[1;31m-> HMAC ERROR: %+v\033[0m\n\n"
)

// CEILINGVALUE : Valor máximo para ser usado como `pModulusValue`
var CEILINGVALUE int = 20

// RAND : Reconfigura Seed
var RAND *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Message : Struct para fazer o parse das mensagens
type Message struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Nonce int    `json:"nonce"`
	HMAC  string `json:"hmac"`
}

//  HandShake : Struct para fazer o parse dos dados do handshake
type HandShake struct {
	Public int `json:"public"`
}

// Check : verifica a existência de erros
func Check(err error, errMessage string) {
	if err != nil {
		fmt.Println(errMessage, " -> ", err.Error())
		os.Exit(1)
	}
}

// GetHMAC ...
func GetHMAC(privateKey string, name string, randString string, nonce int) string {
	hash := hmac.New(sha512.New, []byte(privateKey))
	io.WriteString(hash, fmt.Sprintf("%s.%s.%d", name, randString, nonce))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

func main() {
	fmt.Println("[INFO] Starting Server...")

	var (
		name string
		port int
		alg  string
	)

	flag.StringVar(&name, "name", "client_name", "Nome do Cliente")
	flag.IntVar(&port, "port", 8000, "Porta do Servidor")
	flag.StringVar(&alg, "alg", "diffie-hellman", "Agoritmo para gerar a chave compartilhada")
	flag.Parse()

	fmt.Printf("Launching server at localhost:%d...", port)

	// Inicia servidor na porta 8000
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	Check(err, "Unable to run server!")

	defer ln.Close()
	for {
		// Fica aguardando novos conexões
		conn, err := ln.Accept()
		Check(err, "Unable to accept new conections!")

		// Roda concorrentemente cada nova conexão
		go func(c net.Conn) {
			var dh DiffieHellman
			for {
				msg, err := bufio.NewReader(c).ReadString('\n')
				Check(err, "Unable to get response from client!")

				if msg == "INIT\n" {
					dh.GeneratepModulusValue(CEILINGVALUE)
					dh.GenerategBaseValue()
					dh.GeneratePrivateValue()
					dh.GeneratePublicValue()

					// Envia p, b e o valor público gerado
					c.Write([]byte(fmt.Sprintf(`{"modulus": %d, "base": %d, "public": %d}`+"\n",
						dh.pModulusValue, dh.gBaseValue, dh.publicValue)))
					// Aguarda valor público do cliente
					msg, err := bufio.NewReader(c).ReadString('\n')
					Check(err, "Unable to get response from server!")

					var hs HandShake
					json.Unmarshal([]byte(msg), &hs)

					dh.GenerateSharedPrivateKey(hs.Public)

					fmt.Printf("\nDH: %+v\n", dh)

					fmt.Printf("\nHS: %+v\n", hs)

				} else if msg == "END\n" {
					c.Close()
					return
				} else {
					var message Message
					json.Unmarshal([]byte(msg), &message)
					// output message received

					if GetHMAC(dh.sharedPrivateKey, message.Name, message.Text, message.Nonce) == message.HMAC {
						fmt.Printf(OKCOLOR, message)
						c.Write([]byte("OK" + "\n"))
					} else {
						fmt.Printf(ERRORCOLOR, message)
						c.Write([]byte("HMAC ERROR" + "\n"))
					}
				}
			}
		}(conn)
	}
}
