FROM golang:1.26-alpine

WORKDIR /app

COPY . .
RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "./cmd/app"]
