# climate-check-go

Uma API REST em Go que integra com o [OpenWeatherMap](https://openweathermap.org/api) para fornecer informações meteorológicas em tempo real e ajudar o usuário a planejar sua rotina com base no clima.

## Estrutura do projeto

```
climate-check-go/
├── main.go                             # Ponto de entrada da aplicação
├── config/
│   └── config.go                       # Carregamento de configurações via variáveis de ambiente
└── internal/
    ├── controllers/
    │   └── weather_controller.go       # Handlers HTTP (Gin)
    ├── services/
    │   └── weather_service.go          # Lógica de negócio e sugestões de rotina
    └── repositories/
        └── weather_repository.go       # Acesso à API do OpenWeatherMap
```

## Pré-requisitos

* Go 1.21+
* Chave de API gratuita em <https://openweathermap.org/api>

## Configuração

1. Copie o arquivo de exemplo e preencha com sua chave de API:

```bash
cp .env.example .env
```

2. Edite `.env`:

```env
WEATHER_API_KEY=sua_chave_aqui
```

## Executando

```bash
go run main.go
```

O servidor sobe em `http://localhost:8080` por padrão.

## Endpoints

### `GET /api/v1/weather?city={cidade}`

Retorna o clima atual e uma sugestão de rotina para a cidade informada.

**Exemplo:**

```bash
curl "http://localhost:8080/api/v1/weather?city=São Paulo"
```

**Resposta:**

```json
{
  "city": "São Paulo",
  "country": "BR",
  "temperature": 24.5,
  "feels_like": 25.1,
  "temp_min": 22.0,
  "temp_max": 27.0,
  "humidity": 72,
  "description": "nublado",
  "wind_speed": 3.5,
  "visibility": 10000,
  "suggestion": "Temperatura agradável. Ótimo para atividades ao ar livre."
}
```

---

### `GET /api/v1/weather/forecast?city={cidade}`

Retorna a previsão do tempo para os próximos 5 dias, com um resumo por dia e sugestão de rotina.

**Exemplo:**

```bash
curl "http://localhost:8080/api/v1/weather/forecast?city=Curitiba"
```

**Resposta:**

```json
{
  "city": "Curitiba",
  "country": "BR",
  "days": [
    {
      "date": "2024-07-01",
      "temp_min": 8.0,
      "temp_max": 15.0,
      "humidity": 80,
      "description": "chuva leve",
      "wind_speed": 4.5,
      "suggestion": "Faz frio. Vista-se em camadas e mantenha-se aquecido."
    }
  ]
}
```

## Testes

```bash
go test ./...
```
