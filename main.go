package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/luanorlando/desafio-go-2/server"
)

func handleConcuerncy(cep string) {
	ch := make(chan server.ApiResult)

	go server.RequestBrasilCep(cep, ch)
	go server.RequestViaCep(cep, ch)

	select {
	case result := <-ch:
		fmt.Printf("API mais rápida foi %s, aqui está o resultado: %s\n", result.Source, result.Body)

	case <-time.After(time.Second):
		fmt.Printf("Timeout: Nenhuma das APIS responderam antes de 1 segundo")
	}

	go func() { <-ch }()
}

func handler(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.URL.Path, "/")

	if len(params) != 2 || params[1] == "" {
		http.Error(w, "CEP não encontrado na url", http.StatusBadRequest)
		return
	}

	cep := string(params[1])

	handleConcuerncy(cep)
}

func main() {
	// curl http://localhost:8080/01153000
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
