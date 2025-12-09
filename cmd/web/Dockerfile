FROM golang:1.25.4-alpine AS build

WORKDIR /app

COPY . /app
RUN go build -o oto-web ./cmd/web

EXPOSE 2929
CMD ["./oto-web", "-h", "0.0.0.0", "-p", "2929", "-path", "/app/web"]
