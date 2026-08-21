FROM golang:1.26.4-alpine3.24 AS builder
#imagem da aplicação
WORKDIR /app
#onde vai ser executada as instruções
COPY go.mod go.sum .
#os arquivos que vão ser copiados para a imagem
RUN go mod download
#copia tudo, incluindo o .git (necessário para commit/push de dentro do container)
COPY . .
RUN go build -o /app/mensageforia cmd/server/main.go

FROM alpine:3.24
#imagem final que fica rodando
WORKDIR /app

# git é necessário para o commit/push automático; ca-certificates para HTTPS no push
RUN apk add --no-cache git ca-certificates

COPY --from=builder /app/mensageforia .

# clone "de verdade": leva o .git do build context para a imagem final
COPY --from=builder /app/.git ./.git

# entrypoint configura user.name/email e injeta o GITHUB_TOKEN no remote antes de subir o app
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]
CMD ["./mensageforia"]
