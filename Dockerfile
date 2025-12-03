FROM golang:1.25.4-alpine AS BUILD

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o oto_api ./cmd/api

EXPOSE 9090
CMD ["./oto_api", "-h", "0.0.0.0", "-p", "9090"]
