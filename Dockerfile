FROM golang:1.26.4-alpine3.24 AS builder
#imagem da aplicação
WORKDIR /app
#onde vai ser executada as instruções
COPY go.mod go.sum .
#os arquivos que vão ser copiados para a imagem
RUN go mod download
#ececuta os comando no shell da aplicação
COPY . .
#copia para a raiz
RUN go build cmd/api/main.go
#executa o biniario do go 
FROM alpine:3.24
#a imagem do docker que ficara rodando
WORKDIR /app

COPY --from=builder /app/main . 

ENTRYPOINT ["./main"]

