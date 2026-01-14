package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	exchange, err := fetchExchangeRate()
	if err != nil {
		// Retorna erro 500 caso a busca falhe
		panic(err.Error())
	}

	// Retorna o bid da cotação
	bid := (*exchange)["USDBRL"].Bid

	// Salvar o bid em um arquivo chamado cotacao.txt
	err = os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %v", bid)), 0o644)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Cotação do dólar: %+v\n", bid)
}

func fetchExchangeRate() (*ExchangeResponse, error) {
	// Cria um contexto com timeout de 200 milissegundos
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel() // Garante que o contexto seja cancelado ao final da função

	// URL da API de cotação USD-BRL
	url := "http://localhost:8080/cotacao"

	// Cria a requisição HTTP com o contexto (para respeitar o timeout)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Executa a requisição HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Garante que o body seja fechado ao final

	// Lê todo o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Deserializa o JSON para a struct ExchangeResponse
	var resps ExchangeResponse

	errs := json.Unmarshal(body, &resps)
	if errs != nil {
		return nil, errs
	}

	return &resps, nil
}
