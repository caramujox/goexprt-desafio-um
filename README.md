# goexprt-desafio-um

Olá dev, tudo bem?

Neste desafio vamos aplicar o que aprendemos sobre webserver http, contextos,
banco de dados e manipulação de arquivos com Go.

Você precisará nos entregar dois sistemas em Go:

- client.go
- server.go

Os requisitos para cumprir este desafio são:

O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.

O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida deverá retornar no formato JSON o resultado para o cliente.

Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.

O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.

Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.

O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}

O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080.

Ao finalizar, envie o link do repositório para correção.

---

## 📋 Instruções de Execução

### 🛠 Tecnologias Utilizadas

- Go 1.21+
- GORM (ORM)
- SQLite (banco de dados)
- Context (controle de timeout)
- Net/HTTP (servidor e cliente HTTP)

### 📁 Estrutura do Projeto

```
goexprt-desafio-um/
├── client/
│   └── client.go          # Cliente HTTP
├── server/
│   └── server.go          # Servidor HTTP + API
├── go.mod                 # Dependências
├── go.sum
├── README.md
├── docker-compose.yml
├── exchange.db            # Banco SQLite (gerado automaticamente)
└── cotacao.txt           # Arquivo de cotação (gerado pelo client)
```

### ⚙️ Pré-requisitos

- [Go 1.21+](https://golang.org/dl/)
- Terminal/Prompt de comando

### 🚀 Como Executar

#### 1. Clone o repositório
```bash
git clone <url-do-repositorio>
cd goexprt-desafio-um
```

#### 2. Instale as dependências
```bash
go mod tidy
```

#### 3. (Opcional) Suba o ambiente Docker
Se preferir usar Docker para o SQLite:
```bash
docker-compose up -d
```

#### 4. Execute o servidor (Terminal 1)
```bash
go run server/server.go
```

Você verá a saída:
```
2026/01/02 04:24:46 Starting Database connection
2026/01/02 04:24:46 Database connected successfully
2026/01/02 04:24:46 Initializing HTTP server
2026/01/02 04:24:46 Server Init -> Started at port :8080
```

#### 5. Execute o cliente (Terminal 2)
```bash
go run client/client.go
```

Você verá a saída:
```
2026/01/02 04:25:15 Cotacao received: 6.1234
```

#### 6. (Se usando Docker) Para parar o ambiente
```bash
docker-compose down
```

#### 7. Verifique os resultados

**Arquivo cotacao.txt:**
```
Dólar: 6.1234
```

**Banco de dados:** 
- **Local**: As cotações são persistidas no arquivo `exchange.db`
- **Docker**: As cotações são persistidas no volume `./data/exchange.db`

###  Testando Manualmente

#### Teste direto da API do servidor:
```bash
curl http://localhost:8080/cotacao
```

Resposta esperada:
```json
{"bid":"6.1234"}
```

### ⏱️ Configurações de Timeout

| Componente | Timeout | Descrição |
|------------|---------|-----------|
| Server → API Externa | 200ms | Busca cotação na API AwesomeAPI |
| Server → Banco SQLite | 10ms | Persistência no banco de dados |
| Client → Server | 300ms | Requisição do client ao server |

### 📊 Logs de Timeout

O sistema registra logs específicos quando timeouts são atingidos:

```bash
# Timeout na API externa (200ms)
2026/01/02 04:25:15 API request timeout exceeded (200ms)

# Timeout no banco (10ms)  
2026/01/02 04:25:15 Database operation timeout exceeded (10ms)

# Timeout no client (300ms)
2026/01/02 04:25:15 Client request timeout exceeded (300ms)
```

### 🧪 Checklist de Requisitos

✅ Client realiza requisição HTTP no server solicitando cotação  
✅ Server consome API https://economia.awesomeapi.com.br/json/last/USD-BRL  
✅ Server retorna JSON com resultado para o cliente  
✅ Server registra cotações no banco SQLite com timeout de 10ms  
✅ API externa possui timeout de 200ms  
✅ Client recebe apenas o campo "bid"  
✅ Client possui timeout de 300ms  
✅ Logs de erro para timeouts insuficientes  
✅ Client salva cotação em "cotacao.txt" no formato "Dólar: {valor}"  
✅ Endpoint /cotacao na porta 8080  
✅ Uso de package "context"

---

**Desenvolvido como parte do desafio Go Expert** 🚀
