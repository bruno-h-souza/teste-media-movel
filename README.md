# Teste Média Móvel (MMS)

Este projeto é uma aplicação em Go que consome dados da API do Mercado Bitcoin para calcular a Média Móvel Simples (MMS) de 20, 50 e 200 dias para os pares de criptomoedas **BTC-BRL** e **ETH-BRL**.

A aplicação é dividida em dois modos de execução:
1. **Job de Sincronização**: Busca o histórico de velas (candles) na API do Mercado Bitcoin e salva as médias calculadas no banco de dados.
2. **API REST**: Disponibiliza os dados de indicadores calculados e armazenados através de endpoints HTTP.

## 🛠 Pré-requisitos

- **Go** (versão 1.26 ou superior)
- **Docker** e **Docker Compose** (para subir o banco de dados MySQL)

## ⚙️ Configuração do Ambiente (Infraestrutura)

Antes de rodar a aplicação, você precisa subir a infraestrutura banco de dados utilizando o Docker Compose. 
Na raiz do projeto, execute o comando:

```bash
docker-compose up -d
```
Isso iniciará um container com o **MySQL** (criando o banco de dados e executando os scripts de inicialização).

## 🚀 Como Executar o Projeto

O projeto utiliza a flag `--mode` para definir qual parte do sistema será inicializada.

### 1. Executando o Job de Sincronização
Antes de consultar a API, é recomendável popular o banco de dados local com os dados do Mercado Bitcoin. Execute o Job utilizando:

```bash
go run main.go --mode job
```
*Nota: O Job fará consultas na API do Mercado Bitcoin do último ano até a data atual, calculando e salvando as médias. Isso pode levar alguns minutos devido ao limite de requisições da API externa.*

### 2. Iniciando a API REST
Após popular o banco (ou em paralelo), inicie o servidor web para habilitar os endpoints de consulta:

```bash
go run main.go --mode api
```
A API estará rodando por padrão na porta `8080`.

## 📖 Documentação da API (Swagger)

Com o servidor da API rodando, você pode visualizar e interagir com a documentação dos endpoints construída com o Swagger acessando o link abaixo em seu navegador:

👉 **http://localhost:8080/swagger/index.html**

Lá você encontrará os detalhes de requisição e resposta para o endpoint `/:pair/mms`.

## 🧪 Executando os Testes

O projeto possui testes unitários que cobrem repositórios, serviços (com mocks) e lógicas de sincronização. Para rodar todos os testes e visualizar a saída detalhada, execute:

```bash
go test -v ./...
```