package server

import (
	"fmt"
	"io"
	"net/http"
)

type ApiResult struct {
	Source string
	Body   string
}

func RequestViaCep(cep string, ch chan<- ApiResult) {
	// time.Sleep(time.Second)
	url := "http://viacep.com.br/ws/" + cep + "/json/"

	emptyResult := ApiResult{
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
		ch <- ApiResult{
			Source: "Via CEP",
			Body:   fmt.Sprintf("Erro: status %d recebido da API", resp.StatusCode),
		}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- emptyResult
	}

	ch <- ApiResult{
		Source: "Via CEP",
		Body:   string(body),
	}
}

func RequestBrasilCep(cep string, ch chan<- ApiResult) {
	// time.Sleep(time.Second)
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	req, err := http.NewRequest("GET", url, nil)

	emptyResult := ApiResult{
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
		ch <- ApiResult{
			Source: "Brasil CEP",
			Body:   fmt.Sprintf("Erro: status %d recebido da API", resp.StatusCode),
		}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- emptyResult
	}

	ch <- ApiResult{
		Source: "Brasil CEP",
		Body:   string(body),
	}
}
