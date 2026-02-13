package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type apiResult struct {
	Source string
	Body   string
}

func requestBrasilCep(cep string, ch chan<- apiResult) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	req, err := http.NewRequest("GET", url, nil)

	emptyResult := apiResult{
		Source: "",
		Body:   "",
	}

	if err != nil {
		ch <- emptyResult
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ch <- emptyResult
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- apiResult{
			Source: "Brasil CEP",
			Body:   fmt.Sprintf("Erro: status %d recebido da API", resp.StatusCode),
		}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- emptyResult
	}

	ch <- apiResult{
		Source: "Brasil CEP",
		Body:   string(body),
	}
}

func requestViaCep(cep string, ch chan<- apiResult) {
	url := "http://viacep.com.br/ws/" + cep + "/json/"

	emptyResult := apiResult{
		Source: "",
		Body:   "",
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- emptyResult
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ch <- emptyResult
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- apiResult{
			Source: "Via CEP",
			Body:   fmt.Sprintf("Erro: status %d recebido da API", resp.StatusCode),
		}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- emptyResult
	}

	ch <- apiResult{
		Source: "Via CEP",
		Body:   string(body),
	}
}

func handleConcuerncy(cep string) {
	ch := make(chan apiResult)

	go requestBrasilCep(cep, ch)
	go requestViaCep(cep, ch)

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

	cep := strings.TrimSpace(params[1])

	handleConcuerncy(cep)
}

func main() {
	// http://localhost:8080/cep
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
