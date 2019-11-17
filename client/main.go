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
	"strings"
	"time"
)

// CEILINGVALUE : Valor máximo para ser usado como `pModulusValue`
var CEILINGVALUE int = 20

// RAND : Reconfigura Seed
var RAND *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// LETTERBYTES : Para gerar strings pseudo-aleatórias
const LETTERBYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	// OKCOLOR : Para mostrar OK verde
	OKCOLOR = "[Message %d]: \033[1;32m%s\033[0m"
	// ERRORCOLOR : Para mostrar HMAC ERROR vermelho
	ERRORCOLOR = "[Message %d]: \033[1;31m%s\033[0m"
)

// Message : Struct para fazer o parse das mensagens
type Message struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Nonce int    `json:"nonce"`
	HMAC  string `json:"hmac"`
}

func (m Message) String() string {
	return fmt.Sprintf(
		`{"name" : "%s", "text" : "%s", "nonce" : %d, "hmac" : "%s"}`,
		m.Name, m.Text, m.Nonce, m.HMAC,
	)
}

// HandShake : Struct para fazer o parse dos dados do handshake
type HandShake struct {
	Modulus int `json:"modulus"`
	Base    int `json:"base"`
	Public  int `json:"public"`
}

// GenerateRandString : Gera uma string pseudo-aleatória de tamanho n
func GenerateRandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LETTERBYTES[RAND.Intn(len(LETTERBYTES))]
	}
	return string(b)
}

// GetHMAC : Gera o HMAC
func GetHMAC(privateKey string, name string, randString string, nonce int) string {
	hash := hmac.New(sha512.New, []byte(privateKey))
	io.WriteString(hash, fmt.Sprintf("%s.%s.%d", name, randString, nonce))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

// Check : verifica a existência de erros
func Check(err error, errMessage string) {
	if err != nil {
		fmt.Println(errMessage, " -> ", err.Error())
		os.Exit(1)
	}
}

func main() {
	var (
		name        string
		ip          string
		port        int
		nMessages   int
		lenMessages int
		alg         string
	)

	flag.StringVar(&name, "name", "client_name", "Nome do Cliente")
	flag.StringVar(&ip, "ip", "localhost", "IP do Servidor")
	flag.IntVar(&port, "port", 8000, "Porta do Servidor")
	flag.IntVar(&nMessages, "n_messages", 10, "Quantidade de Mensagens a serem enviadas")
	flag.IntVar(&lenMessages, "len_messages", 10, "Tamanho das Mensagens a serem enviadas")
	flag.StringVar(&alg, "alg", "diffie-hellman", "Agoritmo para gerar a chave compartilhada")
	flag.Parse()

	// Conecta ao Socket
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	Check(err, "Unable to connect to server!")

	var dh DiffieHellman

	// Inicia conexão com o Server enviando `INIT`
	fmt.Fprintf(conn, "INIT"+"\n")
	// Pega resposta
	message, err := bufio.NewReader(conn).ReadString('\n')
	Check(err, "Unable to get response from server!")

	// Pega dados para iniciar o handshake
	var hs HandShake
	json.Unmarshal([]byte(strings.TrimRight(message, "\n")), &hs)

	dh.SetpModulusValue(hs.Modulus)
	dh.SetgBaseValue(hs.Base)
	dh.GeneratePrivateValue()
	dh.GeneratePublicValue()
	dh.GenerateSharedPrivateKey(hs.Public)

	// Testar se a chave privada compartilhada está gerando corretamente
	// fmt.Printf("\nDH: %+v\n", dh)
	// fmt.Printf("\nHS: %+v\n", hs)

	// Envia valor público para o server
	fmt.Fprintf(conn, fmt.Sprintf(`{"public": %d}`+"\n", dh.publicValue))

	for i := 0; i < nMessages; i++ {
		randString := GenerateRandString(lenMessages)
		// TODO: Tornar o Nonce incremental
		nonce := 4242

		msg := &Message{
			Name:  name,
			Text:  randString,
			Nonce: nonce,
			HMAC:  GetHMAC(dh.sharedPrivateKey, name, randString, nonce),
		}
		// Envia mensagem pro Server
		fmt.Fprintf(conn, msg.String()+"\n")

		// Espera resposta
		message, err := bufio.NewReader(conn).ReadString('\n')
		Check(err, "Unable to get response from server!")

		if message == "OK\n" {
			fmt.Printf(OKCOLOR, i+1, message)
		} else {
			fmt.Printf(ERRORCOLOR, i+1, message)
		}

	}
	// Fecha conexão com o server
	fmt.Fprintf(conn, "END"+"\n")
	return
}
