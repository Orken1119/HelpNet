FROM golang:1.23

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Копируем оставшиеся файлы проекта
COPY . /app
# Компилируем приложение
RUN go build -tags netgo -ldflags '-s -w' -o app ./cmd

# Открываем нужный порт (например, 8080)
EXPOSE 1140

# Запускаем приложение
CMD ["./app"]