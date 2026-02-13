package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	fmt.Println(parts)
	if len(parts) != 1 || parts[0] == "" {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
}

func requestBrasilCep(cep string, ch chan<- string) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		ch <- ""
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ch <- ""
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- ""
	}

	ch <- string(body)
}

func requestViaCep(cep string, ch chan<- string) {
	url := "http://viacep.com.br/ws/" + cep + "/json/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- ""
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ch <- ""
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- ""
	}

	ch <- string(body)
}

func handleConcuerncy(cep string) {
	ch := make(chan string)
	ch2 := make(chan string)
	go requestBrasilCep(cep, ch)
	go requestViaCep(cep, ch2)

	select {
	case result := <-ch:
		fmt.Printf("API mais rápida foi Brasil cep, aqui está o resultado: %d", result)
	case result := <-ch2:
		fmt.Printf("API mais rápida foi do Via cep, aqui está o resultado: %d", result)
	case <-time.After(time.Second):
		fmt.Printf("Timeout: Nenhuma das APIS responderam antes de 1 segundo")

	}

}

func handler(w http.ResponseWriter, r *http.Request) {

	params := strings.Split(r.URL.Path, "/")
	fmt.Println(params)

	if len(params) != 2 || params[1] == "" {
		http.Error(w, "CEP não encontrado na url", http.StatusBadRequest)
		return
	}

	cep := params[4]

	handleConcuerncy(cep)
}

func main() {
	// http://localhost:8080/cep
	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":8080", nil)
	handleConcuerncy("15291300")
}
