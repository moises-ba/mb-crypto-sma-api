# MB Crypto SMA API

API em **Go (Golang)** que calcula a **Média Móvel Simples (SMA)** dos preços fechados de criptomoedas (BTC e ETH) com base em dados do Mercado Bitcoin.

A média é calculada para períodos específicos (20, 50 ou 200 dias) a partir de uma data de referência.

---

## Funcionalidades

- Consulta preços fechados diários de criptomoedas (BTC e ETH).
- Calcula a **SMA (Simple Moving Average)** para períodos de 20, 50 ou 200 dias.
- Retorna resultados em **JSON puro**, pronto para consumo em front-end ou integrações.
- Validação de parâmetros de entrada (`days` e `dt_reference`).
- Cache em memória para otimização de performance.
- Estrutura modular com **adapter**, **service**, **strategy** e **controller**.

---

## Estrutura do Projeto

```text
cmd/                # Inicialização da aplicação
internal/
├─ adapter/        # Integração com APIs externas (Mercado Bitcoin)
├─ cache/          # Interface e implementação de cache
├─ controller/     # Handlers HTTP
├─ domain/         # Tipos e entidades do domínio
├─ dto/            # Data Transfer Objects
├─ errors/         # Erros customizados da API
├─ http/           # Cliente HTTP customizado
├─ service/        # Regras de negócio
├─ strategy/       # Estratégias de cálculo (SMA)
└─ log/            # Logging

---

## Como Executar

### Requisitos

- **Docker** e **Docker Compose**
- Go >= 1.20 (para desenvolvimento local, opcional se usar Docker)

### Executando com Docker Compose

```bash
git clone https://github.com/moises-ba/mb-crypto-sma-api
cd mb-crypto-sma-api
docker compose up
```


### chamando a aplicação
ex:
```
curl 'http://localhost:8080/v1/sma/200/2026-01-06'
```