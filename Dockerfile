FROM golang:1.23-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd/server/main.go

EXPOSE 8080

CMD ["/app/main"]