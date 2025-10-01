FROM golang:1.25

WORKDIR /app

COPY . .

RUN go build -o download-service ./cmd/server

EXPOSE 8080

CMD ["./download-service"]