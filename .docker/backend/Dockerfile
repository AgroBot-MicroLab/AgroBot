FROM golang:1.23.1 AS dev
WORKDIR /var/www

COPY backend/go.* ./
RUN go mod download || true
COPY backend/ .
EXPOSE 8080
CMD ["go","run","./cmd/api"]

FROM golang:1.22 AS build
WORKDIR /src
COPY backend/go.* ./
RUN go mod download || true
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/api ./cmd/api

FROM gcr.io/distroless/static:nonroot
WORKDIR /var/www
COPY --from=build /out/api ./api
EXPOSE 8080
USER nonroot
ENTRYPOINT ["./api"]

