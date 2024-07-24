FROM golang:1.22.3

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем исходный код приложения в рабочую директорию контейнера
COPY . .

# Компилируем приложение
RUN GOARCH=amd64 GOOS=linux go build -o socket-service ./cmd/socket-service/main.go

# Указываем команду для запуска приложения
CMD ["./socket-service"]
