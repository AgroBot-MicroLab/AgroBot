FROM golang:1.24.6 AS dev

WORKDIR /var/www

RUN go install github.com/air-verse/air@v1.62.0

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

