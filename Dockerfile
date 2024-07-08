FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go get -d -v ./...
RUN go build -o api ./cmd/api

RUN chmod +x ./entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
