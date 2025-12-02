FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o 2cents cmd/api/main.go

EXPOSE 8080

RUN chmod +x 2cents

CMD ["./2cents"]