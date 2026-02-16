package server

import (
	"fmt"
	"io"
	"net/http"
)

type apiResult struct {
	Source string
	Body   string
}

func RequestViaCep(cep string, ch chan<- apiResult) {
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

func RequestBrasilCep(cep string, ch chan<- apiResult) {
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
