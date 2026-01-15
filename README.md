# Sistema de Cotação do Dólar

Sistema cliente-servidor em Go para consulta e armazenamento da cotação do dólar (USD-BRL).

## 📋 Descrição do Projeto

Este projeto é composto por dois módulos independentes que trabalham em conjunto:

### 🖥️ dolar-server

Servidor HTTP que atua como intermediário entre clientes e a API pública de cotações. Suas responsabilidades incluem:

- Expõe um endpoint REST (`/cotacao`) para consulta da cotação do dólar
- Busca dados em tempo real da API pública [AwesomeAPI](https://economia.awesomeapi.com.br)
- Persiste cada cotação consultada em banco de dados SQLite local
- Implementa timeout de 200ms para requisições à API externa
- Implementa timeout de 10ms para operações de banco de dados

### 💻 dolar-client

Cliente HTTP que consome o servidor e registra a cotação localmente:

- Consulta o endpoint `/cotacao` do servidor local
- Extrai o valor de compra (bid) da cotação
- Salva o resultado em arquivo `cotacao.txt` na raiz do projeto
- Implementa timeout de 300ms para requisições ao servidor

## 🔧 Requisitos

- **Go**: versão 1.25.5 ou superior
- **SQLite**: instalado no sistema (usado pelo dolar-server)

## 📦 Estrutura do Projeto

```
cotacao-dolar/
├── dolar-server/
│   ├── main.go              # Servidor HTTP e lógica de banco de dados
│   ├── exchange_rate.go     # Estruturas de dados (ExchangeRate)
│   ├── go.mod               # Dependências do servidor
│   └── test.db              # Banco SQLite (criado automaticamente)
├── dolar-client/
│   ├── main.go              # Cliente HTTP
│   ├── exchange_rate.go     # Estruturas de dados (ExchangeRate)
│   ├── go.mod               # Dependências do cliente
│   └── cotacao.txt          # Arquivo de saída (criado automaticamente)
└── README.md
```

## 🚀 Instalação e Execução

### Passo 1: Clonar o Repositório (se ainda não tiver)

```bash
git clone <url-do-repositorio>
cd cotacao-dolar
```

### Passo 2: Configurar o dolar-server

```bash
# Navegar até o diretório do servidor
cd dolar-server

# Baixar dependências
go mod download

# Executar o servidor
go run .
```

O servidor iniciará na porta **8080** e ficará aguardando requisições.

**Saída esperada:**
```
(O servidor não exibe mensagem de inicialização, mas estará rodando)
```

### Passo 3: Configurar o dolar-client (em outro terminal)

```bash
# Em um novo terminal, navegar até o diretório do cliente
cd dolar-client

# Baixar dependências (se necessário)
go mod download

# Executar o cliente
go run .
```

**Saída esperada:**
```
Cotação do dólar: 5.8234
```

O arquivo `cotacao.txt` será criado/atualizado com o conteúdo:
```
Dólar: 5.8234
```

## 📡 API do dolar-server

### Endpoints Disponíveis

#### `GET /hello`
Endpoint de verificação do servidor.

**Resposta:**
```
Bem-vindo ao servidor de cotação do dólar!
```

#### `GET /cotacao`
Retorna a cotação atual do dólar em relação ao real.

**Resposta (JSON):**
```json
{
  "USDBRL": {
    "code": "USD",
    "codein": "BRL",
    "name": "Dólar Americano/Real Brasileiro",
    "high": "5.8456",
    "low": "5.7891",
    "varBid": "0.0234",
    "pctChange": "0.4",
    "bid": "5.8234",
    "ask": "5.8267",
    "timestamp": "1705342800",
    "create_date": "2024-01-15 14:20:00"
  }
}
```

**Campos principais:**
- `bid`: Valor de compra (usado pelo cliente)
- `ask`: Valor de venda
- `high`: Máxima do dia
- `low`: Mínima do dia

## 🗄️ Banco de Dados

O servidor utiliza SQLite com GORM. O banco é criado automaticamente no primeiro uso:

**Arquivo:** `dolar-server/test.db`

**Tabela:** `exchange_rates`

Cada requisição bem-sucedida ao endpoint `/cotacao` resulta em um novo registro no banco.

## ⚙️ Configurações e Timeouts

### dolar-server
- **Porta:** 8080
- **Timeout API externa:** 200ms
- **Timeout banco de dados:** 10ms
- **API externa:** `https://economia.awesomeapi.com.br/json/last/USD-BRL`

### dolar-client
- **Servidor alvo:** `http://localhost:8080/cotacao`
- **Timeout requisição:** 300ms
- **Arquivo de saída:** `cotacao.txt`

## 🧪 Testando o Sistema

### 1. Verificar se o servidor está rodando:
```bash
curl http://localhost:8080/hello
```

### 2. Consultar cotação manualmente:
```bash
curl http://localhost:8080/cotacao
```

### 3. Verificar arquivo de saída do cliente:
```bash
cat dolar-client/cotacao.txt
```

### 4. Inspecionar banco de dados (opcional):
```bash
cd dolar-server
sqlite3 test.db "SELECT * FROM exchange_rates ORDER BY id DESC LIMIT 5;"
```

## 🛠️ Compilação (Opcional)

### Compilar o servidor:
```bash
cd dolar-server
go build -o server
./server
```

### Compilar o cliente:
```bash
cd dolar-client
go build -o client
./client
```

## ⚠️ Tratamento de Erros

### Erros comuns:

**1. Servidor não conecta à API externa:**
- Verificar conexão com a internet
- O timeout de 200ms pode ser muito curto em conexões lentas

**2. Cliente não conecta ao servidor:**
- Verificar se o servidor está rodando na porta 8080
- Confirmar que não há firewall bloqueando localhost:8080

**3. Erro ao salvar no banco de dados:**
- Timeout de 10ms muito curto (possível em sistemas lentos)
- Verificar permissões de escrita no diretório

## 📝 Dependências

### dolar-server
```
gorm.io/gorm v1.31.1
gorm.io/driver/sqlite v1.6.0
github.com/mattn/go-sqlite3 v1.14.33
```

### dolar-client
```
(Usa apenas biblioteca padrão do Go)
```

## 📚 Conceitos Demonstrados

Este projeto demonstra:
- ✅ Requisições HTTP com timeout usando `context`
- ✅ Serialização/deserialização JSON
- ✅ Persistência com SQLite e GORM
- ✅ Manipulação de arquivos
- ✅ Estruturação de projetos Go com múltiplos módulos
- ✅ Tratamento de erros
- ✅ Arquitetura cliente-servidor

## 📄 Licença

Este é um projeto educacional para demonstração de conceitos em Go.

---

**Desenvolvido como parte do curso Full Cycle - Go Expert**
