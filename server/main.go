package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

// Body struct to get message from JSON decoder
type Body struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	HMAC    string `json:"hmac"`
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var body Body
		err := decoder.Decode(&body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Header().Set("Connection", "close")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// w.Header().Set("Connection", "close")
		fmt.Printf("%+v\n", body)
	})

	fmt.Println("[INFO] Running server on http://localhost:8000...")
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}
