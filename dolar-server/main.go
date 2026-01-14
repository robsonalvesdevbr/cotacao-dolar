// Package main implementa um servidor HTTP para consulta de cotação do dólar
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// db é a conexão global com o banco de dados SQLite.
// Inicializada uma única vez no início da aplicação.
var db *gorm.DB

// initDatabase inicializa a conexão com o banco de dados
// e executa as migrations necessárias.
func initDatabase() error {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("falha ao conectar ao banco: %w", err)
	}

	// Executa migrations apenas uma vez na inicialização
	return db.AutoMigrate(&ExchangeRate{})
}

// main é o ponto de entrada da aplicação.
// Configura as rotas HTTP e inicia o servidor na porta 8080.
func main() {
	// Inicializa o banco de dados uma única vez
	if err := initDatabase(); err != nil {
		panic(err)
	}

	// Cria um novo multiplexador de rotas HTTP
	mux := http.NewServeMux()

	// Registra os handlers para cada rota
	mux.HandleFunc("/hello", helloWorldHandler)
	mux.HandleFunc("/cotacao", getDollarQuotationHandler)

	// Inicia o servidor HTTP na porta 8080
	http.ListenAndServe(":8080", mux)
}

// helloWorldHandler retorna uma mensagem de boas-vindas.
// É uma rota simples para verificar se o servidor está funcionando.
func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "Bem-vindo ao servidor de cotação do dólar!")
}

// getDollarQuotationHandler é o handler HTTP que retorna a cotação atual do dólar.
// Busca a cotação da API externa e retorna em formato JSON.
func getDollarQuotationHandler(w http.ResponseWriter, r *http.Request) {
	// Busca a cotação do dólar na API externa
	exchange, err := fetchExchangeRate()
	if err != nil {
		// Retorna erro 500 caso a busca falhe
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Define o tipo de conteúdo como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Serializa e envia a resposta em formato JSON
	json.NewEncoder(w).Encode(exchange)
}

// fetchExchangeRate busca a cotação do dólar em relação ao real (USD-BRL)
// na API pública economia.awesomeapi.com.br.
// Utiliza um contexto com timeout de 200ms para evitar requisições longas.
func fetchExchangeRate() (*ExchangeResponse, error) {
	// Cria um contexto com timeout de 200 milissegundos
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel() // Garante que o contexto seja cancelado ao final da função

	// URL da API de cotação USD-BRL
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

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

	// Salva a cotação no banco de dados
	for _, exchange := range resps {
		if err := saveExchangeRate(&exchange); err != nil {
			fmt.Printf("Erro ao salvar cotação: %v\n", err)
		}
	}

	return &resps, nil
}

// saveExchangeRate salva a cotação no banco de dados.
// Utiliza a conexão global já inicializada.
func saveExchangeRate(exchange *ExchangeRate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	result := db.WithContext(ctx).Create(exchange)
	return result.Error
}
