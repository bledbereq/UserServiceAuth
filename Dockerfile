# Этап 1: Сборка
FROM golang:1.22.4-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /usr/local/src

# Устанавливаем необходимые зависимости
RUN apk --no-cache add bash git make task

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной исходный код приложения
COPY . .

# Сборка приложения
RUN go build -o ./bin/app cmd/main.go


# Этап 2: Выполнение
FROM alpine AS runner

# Устанавливаем необходимые зависимости
RUN apk --no-cache add openssl

# Копируем собранный бинарный файл из этапа сборки
COPY --from=builder /usr/local/src/bin/app /app

# Копируем файл конфигурации
COPY config/local.yaml /local.yaml



# Устанавливаем команду для запуска
CMD ["/app", "-config", "/local.yaml"]
