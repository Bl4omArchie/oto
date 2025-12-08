FROM golang:1.25.4-alpine AS BUILD

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o oto-api ./cmd/api

EXPOSE 1515
CMD ["./oto-api", "-h", "0.0.0.0", "-p", "1515"]
