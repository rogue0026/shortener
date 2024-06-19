FROM golang:1.22.4-alpine
RUN apk update && apk add --no-cache git make
WORKDIR /shortener
COPY . .
RUN go mod download
RUN make build
EXPOSE 8080
ENV CONFIG="./configs/local_config.yaml"
ENV PRIVATE_KEY="super secret key for tokens"
CMD ./cmd/bin/shortener

# Выбрать базовый образ
# Установить git
# Создать директорию для размещения проекта
# Скопировать все файлы проекта в созданную директорию
# Выставить порт для работы приложения
# Установить переменные среды для работы приложения
# Выполнить сборку приложения
# Запустить контейнер