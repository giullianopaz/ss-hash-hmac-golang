package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Body struct to get message from JSON decoder
type Body struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	HMAC    string `json:"hmac"`
}

func (body Body) String() string {
	return fmt.Sprintf(
		`{"name" : "%s", "message" : "%s", "hmac" : "%s"}`,
		body.Name, body.Message, body.HMAC,
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

	host := fmt.Sprintf("http://%s:%d", ip, port)

	fmt.Println("n_messages: ", nMessages)

	for i := 0; i < nMessages; i++ {
		// wg.Add(1)
		func() {
			msg := &Body{
				Name:    name,
				Message: "teste",
				HMAC:    "hmac---fwfwlfkmlkwfmewlfml",
			}

			req, err := http.NewRequest("POST", host, bytes.NewBuffer([]byte(msg.String())))
			if err != nil {
				log.Fatal("Error reading request. ", err)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Connection", "close")
			client := &http.Client{Timeout: time.Second * 10}

			resp, err := client.Do(req)
			if err != nil {
				log.Fatal("Error reading response. ", err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Error reading body. ", err)
			}
			fmt.Printf("%s", body)
			// wg.Done()
		}()
	}

	// select {}
}
