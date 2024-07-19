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
RUN go build -o message_service ./cmd/message_service/main.go

# Открываем порт, который будет использоваться приложением
EXPOSE 3004

# Указываем команду для запуска приложения
CMD ["./message_service"]
