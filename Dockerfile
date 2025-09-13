# Stage 1: build
FROM golang:1.25-alpine AS builder

# Dependências do sistema
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copia arquivos do Go e baixa dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código
COPY . .

# Build do binário
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Stage 2: imagem final mínima
# FROM gcr.io/distroless/base-debian10
FROM alpine

WORKDIR /app

# Copia o binário do stage builder
COPY --from=builder /app/app .

# Executa a aplicação
CMD ["./app"]
