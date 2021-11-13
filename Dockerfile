FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN export CGO_ENABLED=0 && go build -o test

ENTRYPOINT ["/app/test"]
