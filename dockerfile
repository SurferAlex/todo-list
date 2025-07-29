FROM golang:1.22-alpine

# Установить PostgreSQL клиент
RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сделать скрипт исполняемым
COPY wait-for-db.sh /wait-for-db.sh
RUN chmod +x /wait-for-db.sh

RUN go build -o main ./cmd/main.go

# Переменные окружения для БД
ENV DB_HOST=localhost
ENV DB_PORT=5433
ENV NAME=test_blog
ENV DB_USER=postgres
ENV DB_PASSWORD=qwe1144EodT5

EXPOSE 8080

# Использовать скрипт ожидания
CMD ["/wait-for-db.sh", "postgres", "./main"]
