# Etapa 1: Construir o binário
FROM golang:1.23 AS builder

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar o código do projeto
COPY . .

# Baixar as dependências
RUN go mod tidy

# Compilar o binário
RUN go build -o main ./cmd/api

# Etapa 2: Criar a imagem final
FROM debian:bookworm-slim

# Instalar dependências necessárias
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Definir o diretório de trabalho
WORKDIR /app

# Copiar o binário gerado na etapa de construção
COPY --from=builder /app/main .

# Copiar outros arquivos necessários (por exemplo, .env)
COPY .env .

# Definir a porta que o serviço irá expor
EXPOSE 3000

# Comando para iniciar o serviço
CMD ["./main"]
