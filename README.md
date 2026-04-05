# 📊 Crypto SMA API

API responsável por calcular a média móvel simples (SMA) dos últimos preços fechados de criptomoedas.

---

## 🚀 Endpoint

### 🔹 Listar médias dos últimos preços fechados

Retorna a média dos preços fechados com base na quantidade de dias e uma data de referência.

---

## 📥 Parâmetros

| Parâmetro      | Tipo   | Obrigatório | Descrição                                      |
|----------------|--------|------------|-----------------------------------------------|
| `days`         | int    | ✅         | Quantidade de dias (ex: 20, 50, 200)          |
| `dt_reference` | string | ✅         | Data de referência no formato `YYYY-MM-DD`    |

---

## 📤 Respostas

### ✅ 200 OK

Retorna a média dos preços fechados:

```json
{
  "last_closed_prices_avg": [
    {
      "digital_coin": "BTC",
      "value": "43000.25",
      "days": 20,
      "avg_type": "SMA",
      "reference_date": "2024-01-01"
    },
        {
      "digital_coin": "ETH",
      "value": "23000.25",
      "days": 20,
      "avg_type": "SMA",
      "reference_date": "2024-01-01"
    }
  ]
}

## 🚀 Como Executar o Projeto

Para rodar o projeto localmente, você precisa ter o **Docker** e o **Docker Compose** instalados.  

### Passos:

1. Clone o repositório:

```bash
git clone https://github.com/moises-ba/mb-crypto-sma-api -b develop
cd mb-crypto-sma-api
docker compose up
```