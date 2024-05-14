FROM golang:latest
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/app/
ENTRYPOINT ["/app/main"]
